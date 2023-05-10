package evm

import (
	"errors"
	"strings"

	"github.com/fxamacker/cbor/v2"
)

// TODO: move to github.com/smartcontractkit/ocr2keepers/pkg/types

// LogTriggerUpkeepConfig holds the settings for a log upkeep
type LogTriggerUpkeepConfig struct {
	// Address is required, contract address w/ 0x prefix
	Address string `cbor:"a"`
	// Topic is required, 32 bytes, w/o 0x prefixed
	Topic string `cbor:"sig"`
	// Filter1 is optional, needs to be left-padded to 32 bytes
	Filter1 string `cbor:"f1,omitempty"`
	// Filter2 is optional, needs to be left-padded to 32 bytes
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
