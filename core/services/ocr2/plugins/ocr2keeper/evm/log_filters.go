package evm

import (
	"errors"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

var (
	logRetention = time.Hour // TODO: check what value to use
)

// toLogFilter converts the given upkeep config into a log filter
func toLogFilter(ltcfg *LogTriggerUpkeepConfig, upkeepID string) logpoller.Filter {
	sigs := []common.Hash{
		common.BytesToHash([]byte(ltcfg.Topic)),
	}
	sigs = ltcfg.addFilters(sigs)
	return logpoller.Filter{
		Name:      filterName(upkeepID),
		EventSigs: sigs,
		Addresses: []common.Address{common.BytesToAddress([]byte(ltcfg.Address))},
		Retention: logRetention,
	}
}

func filterName(upkeepID string) string {
	return logpoller.FilterName(upkeepID)
}

// logFiltersProvider manages log filters for upkeeps
type logFiltersProvider struct {
	poller logpoller.LogPoller
}

func newLogFiltersProvider(poller logpoller.LogPoller) *logFiltersProvider {
	return &logFiltersProvider{
		poller: poller,
	}
}

// Register takes an upkeep and register the corresponding filter if applicable
func (lfp *logFiltersProvider) Register(upkeepID string, rawCfg []byte) error {
	cfg := &LogTriggerUpkeepConfig{}
	if err := cfg.Decode(rawCfg); err != nil {
		return err
	}
	if err := cfg.Validate(); err != nil {
		return err
	}
	// remove old filter (if exist) for this upkeep TODO: check if needed
	// _ = lfp.poller.UnregisterFilter(filterName(upkeepID), nil)
	return lfp.poller.RegisterFilter(toLogFilter(cfg, upkeepID))
}

func (lfp *logFiltersProvider) UnRegister(upkeepID string) error {
	return lfp.poller.UnregisterFilter(filterName(upkeepID), nil)
}

// TODO: move?

// LogTriggerUpkeepConfig holds the settings for a log upkeep
type LogTriggerUpkeepConfig struct {
	// Address is a 0x prefixed contract address (required)
	Address string `cbor:"a"`
	// Topic 32 bytes, not 0x prefixed (required)
	Topic string `cbor:"sig"`
	// Filter1 needs to be left-padded to 32 bytes (optional)
	Filter1 string `cbor:"f1,omitempty"`
	// Filter2 needs to be left-padded to 32 bytes (optional)
	Filter2 string `cbor:"f2,omitempty"`
	// Filter3 is wildcard filter if missing
	Filter3 string `cbor:"f3,omitempty"`
}

var (
	ErrContractAddrNotFound = errors.New("invalid contract address: not found")
	ErrContractAddrNoPrefix = errors.New("invalid contract address: not prefixed with 0x")
	ErrTopicIncorrectSize   = errors.New("invalid topic: incorrect size")
	ErrTopicPrefix          = errors.New("invalid topic: prefixed with 0x")
)

func (ltcfg *LogTriggerUpkeepConfig) Validate() error {
	if len(ltcfg.Address) == 0 {
		return ErrContractAddrNotFound
	}
	if strings.Index(ltcfg.Address, "0x") != 0 {
		// ltcfg.Address = fmt.Sprintf("0x%s", ltcfg.Address)
		return ErrContractAddrNoPrefix
	}
	if n := len(ltcfg.Topic); n != 32 {
		return ErrTopicIncorrectSize
	}
	if strings.Index(ltcfg.Topic, "0x") == 0 {
		return ErrTopicPrefix
	}

	// TODO: validate filters

	return nil
}

func (ltcfg *LogTriggerUpkeepConfig) Encode() ([]byte, error) {
	return cbor.Marshal(ltcfg)
}

func (ltcfg *LogTriggerUpkeepConfig) Decode(raw []byte) error {
	return cbor.Unmarshal(raw, ltcfg)
}

func (ltcfg *LogTriggerUpkeepConfig) addFilters(sigs []common.Hash) []common.Hash {
	if len(ltcfg.Filter1) > 0 {
		sigs = append(sigs, common.BytesToHash([]byte(ltcfg.Filter1)))
	}
	return sigs
}
