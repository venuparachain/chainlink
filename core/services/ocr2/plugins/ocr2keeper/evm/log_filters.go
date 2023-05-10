package evm

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

const (
	// logRetention is the amount of time to retain logs for
	logRetention = time.Minute * 5
)

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
	return lfp.poller.RegisterFilter(newLogFilter(cfg, upkeepID))
}

func (lfp *logFiltersProvider) UnRegister(upkeepID string) error {
	return lfp.poller.UnregisterFilter(filterName(upkeepID), nil)
}

// newLogFilter creates logpoller.Filter from the given upkeep config
// NOTE: assuming that the upkeep config is valid
func newLogFilter(ltcfg *LogTriggerUpkeepConfig, upkeepID string) logpoller.Filter {
	sigs := []common.Hash{
		common.BytesToHash([]byte(ltcfg.Topic)),
	}
	sigs = addFilters(sigs, ltcfg)
	return logpoller.Filter{
		Name:      filterName(upkeepID),
		EventSigs: sigs,
		Addresses: []common.Address{common.BytesToAddress([]byte(ltcfg.Address))},
		Retention: logRetention,
	}
}

func addFilters(sigs []common.Hash, ltcfg *LogTriggerUpkeepConfig) []common.Hash {
	if len(ltcfg.Filter1) > 0 {
		sigs = append(sigs, common.BytesToHash([]byte(ltcfg.Filter1)))
	}
	if len(ltcfg.Filter2) > 0 {
		sigs = append(sigs, common.BytesToHash([]byte(ltcfg.Filter2)))
	}
	if len(ltcfg.Filter3) > 0 {
		sigs = append(sigs, common.BytesToHash([]byte(ltcfg.Filter3)))
	}
	return sigs
}

func filterName(upkeepID string) string {
	return logpoller.FilterName(upkeepID)
}
