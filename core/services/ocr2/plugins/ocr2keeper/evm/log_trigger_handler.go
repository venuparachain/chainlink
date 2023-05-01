package evm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"
)

// TODO: move

type LogTriggerUpkeepConfig struct {
	Address string `cbor:"a"`            // required, 0x prefixed address
	Topic   string `cbor:"sig"`          // required, 32 bytes, not 0x prefixed
	Filter1 string `cbor:"f1,omitempty"` // optional, needs to be left-padded to 32 bytes
	Filter2 string `cbor:"f2,omitempty"` // optional, needs to be left-padded to 32 bytes
	Filter3 string `cbor:"f3,omitempty"` // f3 is wild if missing
}

func (ltcfg *LogTriggerUpkeepConfig) Validate() error {
	if len(ltcfg.Address) == 0 {
		return errors.New("invalid contract address: missing")
	}
	if strings.Index(ltcfg.Address, "0x") != 0 {
		// ltcfg.Address = fmt.Sprintf("0x%s", ltcfg.Address)
		return errors.New("invalid contract address: not prefixed with 0x")
	}
	if n := len(ltcfg.Topic); n != 32 {
		return fmt.Errorf("invalid topic: incorrect size %d", n)
	}
	if strings.Index(ltcfg.Topic, "0x") == 0 {
		return fmt.Errorf("invalid topic: prefixed with 0x")
	}

	// TODO: filters

	return nil
}

func (ltcfg *LogTriggerUpkeepConfig) Encode() ([]byte, error) {
	return cbor.Marshal(ltcfg)
}

func (ltcfg *LogTriggerUpkeepConfig) Decode(raw []byte) error {
	return cbor.Unmarshal(raw, ltcfg)
}

func (ltcfg *LogTriggerUpkeepConfig) ToLogFilter(uid types.UpkeepIdentifier) logpoller.Filter {
	sigs := []common.Hash{
		common.BytesToHash([]byte(ltcfg.Topic)),
	}
	sigs = ltcfg.addFilters(sigs)
	return logpoller.Filter{
		Name:      string(uid), // TODO: format?
		EventSigs: sigs,
		Addresses: []common.Address{common.BytesToAddress([]byte(ltcfg.Address))},
	}
}

func (ltcfg *LogTriggerUpkeepConfig) addFilters(sigs []common.Hash) []common.Hash {
	if len(ltcfg.Filter1) > 0 {
		sigs = append(sigs, common.BytesToHash([]byte(ltcfg.Filter1)))
	}
	return sigs
}

type logTriggerHandler struct {
	poller logpoller.LogPoller
}

var _ UpkeepStateHandler = &logTriggerHandler{}

// Handle is invoked upon an upkeep state change, it filters log trigger upkeeps
// and for each upkeep it updates the log filter.
func (handler *logTriggerHandler) Handle(st stateEventType, abilog generated.AbigenLog) error {
	cfg := new(LogTriggerUpkeepConfig)
	uid, err := decodeUpkeepStateEvent(cfg, abilog)
	if err != nil {
		return err
	}
	if err := cfg.Validate(); err != nil {
		return err
	}
	switch st {
	case stateUpkeepCanceled, stateUpkeepPaused:
		return handler.poller.UnregisterFilter(string(uid), nil) // TODO: queryer
	case stateUpkeepRegistered, stateUpkeepUnpaused:
		return handler.poller.RegisterFilter(cfg.ToLogFilter(uid))
	}
	return nil
}
