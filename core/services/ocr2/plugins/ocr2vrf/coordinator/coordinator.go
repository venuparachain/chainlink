package coordinator

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ocr2vrf/dkg"
	ocr2recoverytypes "github.com/smartcontractkit/ocr2vrf/types"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	dkg_wrapper "github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/dkg"
	recovery_beacon "github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/recovery_beacon"
	"github.com/smartcontractkit/chainlink/core/logger"
)

var _ ocr2recoverytypes.CoordinatorInterface = &coordinator{}

var (
	dkgABI            = evmtypes.MustGetABI(dkg_wrapper.DKGMetaData.ABI)
	recoveryBeaconABI = evmtypes.MustGetABI(recovery_beacon.RecoveryBeaconMetaData.ABI)
)

const (
	configSetEvent = "ConfigSet"
)

type coordinator struct {
	lggr logger.Logger

	lp logpoller.LogPoller
	topics
	lookbackBlocks int64
	finalityDepth  uint32

	beaconAddress common.Address

	// We need to keep track of DKG ConfigSet events as well.
	dkgAddress common.Address

	evmClient evmclient.Client
}

// New creates a new CoordinatorInterface implementor.
func New(
	lggr logger.Logger,
	beaconAddress common.Address,
	dkgAddress common.Address,
	client evmclient.Client,
	logPoller logpoller.LogPoller,
) (ocr2recoverytypes.CoordinatorInterface, error) {

	t := newTopics()

	// Add log filters for the log poller so that it can poll and find the logs that
	// we need.
	_, err := logPoller.RegisterFilter(logpoller.Filter{
		EventSigs: []common.Hash{
			t.configSetTopic,
		}, Addresses: []common.Address{beaconAddress, dkgAddress}})
	if err != nil {
		return nil, err
	}

	return &coordinator{
		beaconAddress: beaconAddress,
		dkgAddress:    dkgAddress,
		lp:            logPoller,
		topics:        t,
		evmClient:     client,
		lggr:          lggr.Named("OCR2RecoveryCoordinator"),
	}, nil
}

func (c *coordinator) ContractID(ctx context.Context) common.Address {
	return c.beaconAddress
}

// DKGRecoveryCommittees returns the addresses of the signers and transmitters
// for the DKG and Recovery OCR committees. On ethereum, these can be retrieved
// from the most recent ConfigSet events for each contract.
func (c *coordinator) DKGRecoveryCommittees(ctx context.Context) (dkgCommittee, vrfCommittee ocr2recoverytypes.OCRCommittee, err error) {

	latestRecovery, err := c.lp.LatestLogByEventSigWithConfs(
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

	var recoveryConfigSetLog recovery_beacon.RecoveryBeaconConfigSet
	err = recoveryBeaconABI.UnpackIntoInterface(&recoveryConfigSetLog, configSetEvent, latestRecovery.Data)
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
	for i := range recoveryConfigSetLog.Signers {
		vrfCommittee.Signers = append(vrfCommittee.Signers, recoveryConfigSetLog.Signers[i])
		vrfCommittee.Transmitters = append(vrfCommittee.Transmitters, recoveryConfigSetLog.Transmitters[i])
	}

	for i := range dkgConfigSetLog.Signers {
		dkgCommittee.Signers = append(dkgCommittee.Signers, dkgConfigSetLog.Signers[i])
		dkgCommittee.Transmitters = append(dkgCommittee.Transmitters, dkgConfigSetLog.Transmitters[i])
	}

	return
}

func (c *coordinator) ProvingKeyHash(ctx context.Context) (common.Hash, error) {
	beacon, err := recovery_beacon.NewRecoveryBeacon(c.beaconAddress, c.evmClient)
	if err != nil {
		return [32]byte{}, err
	}

	h, err := beacon.SProvingKeyHash(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "get proving block hash")
	}

	return h, nil
}

func (c *coordinator) KeyID(ctx context.Context) (dkg.KeyID, error) {
	beacon, err := recovery_beacon.NewRecoveryBeacon(c.beaconAddress, c.evmClient)
	if err != nil {
		return [32]byte{}, err
	}

	keyID, err := beacon.SKeyID(&bind.CallOpts{Context: ctx})
	if err != nil {
		return dkg.KeyID{}, errors.Wrap(err, "could not get key ID")
	}
	return keyID, nil
}
