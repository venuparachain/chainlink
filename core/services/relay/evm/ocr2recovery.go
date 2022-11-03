package evm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/sqlx"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/dkg/config"
)

// DKGProvider provides all components needed for a DKG plugin.
type DKGProvider interface {
	relaytypes.Plugin
}

// OCR2RecoveryProvider provides all components needed for a OCR2Recovery plugin.
type OCR2RecoveryProvider interface {
	relaytypes.Plugin
}

// OCR2RecoveryRelayer contains the relayer and instantiating functions for OCR2Recovery providers.
type OCR2RecoveryRelayer interface {
	NewDKGProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (DKGProvider, error)
	NewOCR2RecoveryProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (OCR2RecoveryProvider, error)
}

var (
	_ OCR2RecoveryRelayer  = (*ocr2recoveryRelayer)(nil)
	_ DKGProvider          = (*dkgProvider)(nil)
	_ OCR2RecoveryProvider = (*ocr2recoveryProvider)(nil)
)

// Relayer with added DKG and OCR2Recovery provider functions.
type ocr2recoveryRelayer struct {
	db    *sqlx.DB
	chain evm.Chain
	lggr  logger.Logger
}

func NewOCR2RecoveryRelayer(db *sqlx.DB, chain evm.Chain, lggr logger.Logger) OCR2RecoveryRelayer {
	return &ocr2recoveryRelayer{
		db:    db,
		chain: chain,
		lggr:  lggr,
	}
}

func (r *ocr2recoveryRelayer) NewDKGProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (DKGProvider, error) {
	configWatcher, err := newOCR2RecoveryConfigProvider(r.lggr, r.chain, rargs)
	if err != nil {
		return nil, err
	}
	contractTransmitter, err := newContractTransmitter(r.lggr, rargs, pargs.TransmitterID, configWatcher)
	if err != nil {
		return nil, err
	}

	var pluginConfig config.PluginConfig
	err = json.Unmarshal(pargs.PluginConfig, &pluginConfig)
	if err != nil {
		return nil, err
	}

	return &dkgProvider{
		configWatcher:       configWatcher,
		contractTransmitter: contractTransmitter,
		pluginConfig:        pluginConfig,
	}, nil
}

func (r *ocr2recoveryRelayer) NewOCR2RecoveryProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (OCR2RecoveryProvider, error) {
	configWatcher, err := newOCR2RecoveryConfigProvider(r.lggr, r.chain, rargs)
	if err != nil {
		return nil, err
	}
	contractTransmitter, err := newContractTransmitter(r.lggr, rargs, pargs.TransmitterID, configWatcher)
	if err != nil {
		return nil, err
	}
	return &ocr2recoveryProvider{
		configWatcher:       configWatcher,
		contractTransmitter: contractTransmitter,
	}, nil
}

type dkgProvider struct {
	*configWatcher
	contractTransmitter *ContractTransmitter
	pluginConfig        config.PluginConfig
}

func (c *dkgProvider) ContractTransmitter() types.ContractTransmitter {
	return c.contractTransmitter
}

type ocr2recoveryProvider struct {
	*configWatcher
	contractTransmitter *ContractTransmitter
}

func (c *ocr2recoveryProvider) ContractTransmitter() types.ContractTransmitter {
	return c.contractTransmitter
}

func newOCR2RecoveryConfigProvider(lggr logger.Logger, chain evm.Chain, rargs relaytypes.RelayArgs) (*configWatcher, error) {
	var relayConfig RelayConfig
	err := json.Unmarshal(rargs.RelayConfig, &relayConfig)
	if err != nil {
		return nil, err
	}
	if !common.IsHexAddress(rargs.ContractID) {
		return nil, fmt.Errorf("invalid contract address '%s'", rargs.ContractID)
	}

	contractAddress := common.HexToAddress(rargs.ContractID)
	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
	if err != nil {
		return nil, errors.Wrap(err, "could not get OCR2Aggregator ABI JSON")
	}
	configPoller, err := NewConfigPoller(
		lggr.With("contractID", rargs.ContractID),
		chain.LogPoller(),
		contractAddress)
	if err != nil {
		return nil, err
	}

	offchainConfigDigester := evmutil.EVMOffchainConfigDigester{
		ChainID:         chain.Config().ChainID().Uint64(),
		ContractAddress: contractAddress,
	}

	return newConfigWatcher(
		lggr,
		contractAddress,
		contractABI,
		offchainConfigDigester,
		configPoller,
		chain,
		relayConfig.FromBlock,
		rargs.New,
	), nil
}
