package coordinator

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ocr2vrf/dkg"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	dkg_wrapper "github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/dkg"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"
	"github.com/smartcontractkit/chainlink/core/logger"
)

var _ ocr2vrftypes.CoordinatorInterface = &coordinator{}

var (
	dkgABI            = evmtypes.MustGetABI(dkg_wrapper.DKGMetaData.ABI)
	vrfBeaconABI      = evmtypes.MustGetABI(vrf_beacon.VRFBeaconMetaData.ABI)
	vrfCoordinatorABI = evmtypes.MustGetABI(vrf_coordinator.VRFCoordinatorMetaData.ABI)
)

const (
	// VRF-only events.
	randomnessRequestedEvent            string = "RandomnessRequested"
	randomnessFulfillmentRequestedEvent string = "RandomnessFulfillmentRequested"
	randomWordsFulfilledEvent           string = "RandomWordsFulfilled"
	newTransmissionEvent                string = "NewTransmission"
	outputsServedEvent                  string = "OutputsServed"

	// Both VRF and DKG contracts emit this, it's an OCR event.
	configSetEvent = "ConfigSet"

	// TODO: Add these defaults to the off-chain config, and get better gas estimates
	// (these current values are very conservative).
	cacheEvictionWindowSeconds = 60               // the maximum duration (in seconds) that an item stays in the cache
	batchGasLimit              = int64(5_000_000) // maximum gas limit of a report
	blockGasOverhead           = int64(50_000)    // cost of posting a randomness seed for a block on-chain
	coordinatorOverhead        = int64(50_000)    // overhead costs of running the transmit transaction
)

// block is used to key into a set that tracks beacon blocks.
type block struct {
	blockNumber uint64
	confDelay   uint32
}

type blockInReport struct {
	block
	recentBlockHeight uint64
	recentBlockHash   common.Hash
}

type callback struct {
	blockNumber uint64
	requestID   uint64
}

type callbackInReport struct {
	callback
	recentBlockHeight uint64
	recentBlockHash   common.Hash
}

type coordinator struct {
	lggr logger.Logger

	lp logpoller.LogPoller
	topics
	lookbackBlocks int64
	finalityDepth  uint32

	onchainRouter      VRFBeaconCoordinator
	coordinatorAddress common.Address
	beaconAddress      common.Address

	// We need to keep track of DKG ConfigSet events as well.
	dkgAddress common.Address

	evmClient evmclient.Client

	// set of blocks that have been scheduled for transmission.
	toBeTransmittedBlocks *ocrCache[blockInReport]
	// set of request id's that have been scheduled for transmission.
	toBeTransmittedCallbacks *ocrCache[callbackInReport]
}

// New creates a new CoordinatorInterface implementor.
func New(
	lggr logger.Logger,
	beaconAddress common.Address,
	coordinatorAddress common.Address,
	dkgAddress common.Address,
	client evmclient.Client,
	lookbackBlocks int64,
	logPoller logpoller.LogPoller,
	finalityDepth uint32,
) (ocr2vrftypes.CoordinatorInterface, error) {
	onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, client)
	if err != nil {
		return nil, errors.Wrap(err, "onchain router creation")
	}

	t := newTopics()

	// Add log filters for the log poller so that it can poll and find the logs that
	// we need.
	_, err = logPoller.RegisterFilter(logpoller.Filter{
		EventSigs: []common.Hash{
			t.randomnessRequestedTopic,
			t.randomnessFulfillmentRequestedTopic,
			t.randomWordsFulfilledTopic,
			t.configSetTopic,
			t.outputsServedTopic}, Addresses: []common.Address{beaconAddress, coordinatorAddress, dkgAddress}})
	if err != nil {
		return nil, err
	}

	cacheEvictionWindow := time.Duration(cacheEvictionWindowSeconds * time.Second)

	return &coordinator{
		onchainRouter:            onchainRouter,
		coordinatorAddress:       coordinatorAddress,
		beaconAddress:            beaconAddress,
		dkgAddress:               dkgAddress,
		lp:                       logPoller,
		topics:                   t,
		lookbackBlocks:           lookbackBlocks,
		finalityDepth:            finalityDepth,
		evmClient:                client,
		lggr:                     lggr.Named("OCR2VRFCoordinator"),
		toBeTransmittedBlocks:    NewBlockCache[blockInReport](cacheEvictionWindow),
		toBeTransmittedCallbacks: NewBlockCache[callbackInReport](cacheEvictionWindow),
	}, nil
}

func (c *coordinator) ContractID(ctx context.Context) common.Address {
	return c.beaconAddress
}

// DKGVRFCommittees returns the addresses of the signers and transmitters
// for the DKG and VRF OCR committees. On ethereum, these can be retrieved
// from the most recent ConfigSet events for each contract.
func (c *coordinator) DKGVRFCommittees(ctx context.Context) (dkgCommittee, vrfCommittee ocr2vrftypes.OCRCommittee, err error) {

	latestVRF, err := c.lp.LatestLogByEventSigWithConfs(
		c.configSetTopic,
		c.beaconAddress,
		int(c.finalityDepth),
	)
	if err != nil {
		err = errors.Wrap(err, "latest vrf ConfigSet by sig with confs")
		return
	}

	latestDKG, err := c.lp.LatestLogByEventSigWithConfs(
		c.configSetTopic,
		c.dkgAddress,
		int(c.finalityDepth),
	)
	if err != nil {
		err = errors.Wrap(err, "latest dkg ConfigSet by sig with confs")
		return
	}

	var vrfConfigSetLog vrf_beacon.VRFBeaconConfigSet
	err = vrfBeaconABI.UnpackIntoInterface(&vrfConfigSetLog, configSetEvent, latestVRF.Data)
	if err != nil {
		err = errors.Wrap(err, "unpack vrf ConfigSet into interface")
		return
	}

	var dkgConfigSetLog dkg_wrapper.DKGConfigSet
	err = dkgABI.UnpackIntoInterface(&dkgConfigSetLog, configSetEvent, latestDKG.Data)
	if err != nil {
		err = errors.Wrap(err, "unpack dkg ConfigSet into interface")
		return
	}

	// len(signers) == len(transmitters), this is guaranteed by libocr.
	for i := range vrfConfigSetLog.Signers {
		vrfCommittee.Signers = append(vrfCommittee.Signers, vrfConfigSetLog.Signers[i])
		vrfCommittee.Transmitters = append(vrfCommittee.Transmitters, vrfConfigSetLog.Transmitters[i])
	}

	for i := range dkgConfigSetLog.Signers {
		dkgCommittee.Signers = append(dkgCommittee.Signers, dkgConfigSetLog.Signers[i])
		dkgCommittee.Transmitters = append(dkgCommittee.Transmitters, dkgConfigSetLog.Transmitters[i])
	}

	return
}

// ProvingKeyHash returns the VRF current proving block, in view of the local
// node. On ethereum this can be retrieved from the VRF contract's attribute
// s_provingKeyHash
func (c *coordinator) ProvingKeyHash(ctx context.Context) (common.Hash, error) {
	h, err := c.onchainRouter.SProvingKeyHash(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "get proving block hash")
	}

	return h, nil
}

// KeyID returns the key ID from coordinator's contract
func (c *coordinator) KeyID(ctx context.Context) (dkg.KeyID, error) {
	keyID, err := c.onchainRouter.SKeyID(&bind.CallOpts{Context: ctx})
	if err != nil {
		return dkg.KeyID{}, errors.Wrap(err, "could not get key ID")
	}
	return keyID, nil
}
