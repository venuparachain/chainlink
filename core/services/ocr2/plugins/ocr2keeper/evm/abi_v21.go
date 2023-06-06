package evm

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"

	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
)

type evmRegistryPackerV21 struct {
	abi abi.ABI
}

func NewEvmRegistryPackerV21(abi abi.ABI) *evmRegistryPackerV21 {
	return &evmRegistryPackerV21{abi: abi}
}

// TODO: adjust to 2.1
func (rp *evmRegistryPackerV21) UnpackCheckResult(key ocr2keepers.UpkeepKey, raw string) (EVMAutomationUpkeepResult20, error) {
	var (
		result EVMAutomationUpkeepResult20
	)

	b, err := hexutil.Decode(raw)
	if err != nil {
		return result, err
	}

	out, err := rp.abi.Methods["checkUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		return result, fmt.Errorf("%w: unpack checkUpkeep return: %s", err, raw)
	}

	block, id, err := splitKey(key)
	if err != nil {
		return result, err
	}

	result = EVMAutomationUpkeepResult20{
		Block:    uint32(block.Uint64()),
		ID:       id,
		Eligible: true,
	}

	upkeepNeeded := *abi.ConvertType(out[0], new(bool)).(*bool)
	rawPerformData := *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	result.FailureReason = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	result.GasUsed = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	result.FastGasWei = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	result.LinkNative = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	if !upkeepNeeded {
		result.Eligible = false
	}
	// if NONE we expect the perform data. if TARGET_CHECK_REVERTED we will have the error data in the perform data used for off chain lookup
	if result.FailureReason == UPKEEP_FAILURE_REASON_NONE || (result.FailureReason == UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED && len(rawPerformData) > 0) {
		var ret0 = new(performDataWrapper)
		err = pdataABI.UnpackIntoInterface(ret0, "check", rawPerformData)
		if err != nil {
			return result, err
		}

		result.CheckBlockNumber = ret0.Result.CheckBlockNumber
		result.CheckBlockHash = ret0.Result.CheckBlockhash
		result.PerformData = ret0.Result.PerformData
	}

	// This is a default placeholder which is used since we do not get the execute gas
	// from checkUpkeep result. This field is overwritten later from the execute gas
	// we have for an upkeep in memory. TODO (AUTO-1482): Refactor this
	result.ExecuteGas = 5_000_000

	return result, nil
}

func (rp *evmRegistryPackerV21) UnpackMercuryLookupResult(callbackResp []byte) (bool, []byte, error) {
	typBytes, err := abi.NewType("bytes", "", nil)
	if err != nil {
		return false, nil, fmt.Errorf("abi new bytes type error: %w", err)
	}
	boolTyp, err := abi.NewType("bool", "", nil)
	if err != nil {
		return false, nil, fmt.Errorf("abi new bool type error: %w", err)
	}
	callbackOutput := abi.Arguments{
		{Name: "upkeepNeeded", Type: boolTyp},
		{Name: "performData", Type: typBytes},
	}
	unpack, err := callbackOutput.Unpack(callbackResp)
	if err != nil {
		return false, nil, fmt.Errorf("callback output unpack error: %w", err)
	}

	upkeepNeeded := *abi.ConvertType(unpack[0], new(bool)).(*bool)
	if !upkeepNeeded {
		return false, nil, nil
	}
	performData := *abi.ConvertType(unpack[1], new([]byte)).(*[]byte)
	return true, performData, nil
}

func (rp *evmRegistryPackerV21) UnpackPerformResult(raw string) (bool, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		return false, err
	}

	out, err := rp.abi.Methods["simulatePerformUpkeep"].
		Outputs.UnpackValues(b)
	if err != nil {
		return false, fmt.Errorf("%w: unpack simulatePerformUpkeep return: %s", err, raw)
	}

	return *abi.ConvertType(out[0], new(bool)).(*bool), nil
}

func (rp *evmRegistryPackerV21) UnpackUpkeepInfo(id *big.Int, raw string) (iregistry21.UpkeepInfo, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		return iregistry21.UpkeepInfo{}, err
	}

	out, err := rp.abi.Methods["getUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		return iregistry21.UpkeepInfo{}, fmt.Errorf("%w: unpack getUpkeep return: %s", err, raw)
	}

	info := *abi.ConvertType(out[0], new(iregistry21.UpkeepInfo)).(*iregistry21.UpkeepInfo)

	return info, nil
}

// TODO: adjust to 2.1 if needed
func (rp *evmRegistryPackerV21) UnpackTransmitTxInput(raw []byte) ([]ocr2keepers.UpkeepResult, error) {
	var (
		enc     = EVMAutomationEncoder20{}
		decoded []ocr2keepers.UpkeepResult
		out     []interface{}
		err     error
		b       []byte
		ok      bool
	)

	if out, err = rp.abi.Methods["transmit"].Inputs.UnpackValues(raw); err != nil {
		return nil, fmt.Errorf("%w: unpack TransmitTxInput return: %s", err, raw)
	}

	if len(out) < 2 {
		return nil, fmt.Errorf("invalid unpacking of TransmitTxInput in %s", raw)
	}

	if b, ok = out[1].([]byte); !ok {
		return nil, fmt.Errorf("unexpected value type in transaction")
	}

	if decoded, err = enc.DecodeReport(b); err != nil {
		return nil, fmt.Errorf("error during decoding report while unpacking TransmitTxInput: %w", err)
	}

	return decoded, nil
}

// UnpackLogTriggerConfig unpacks the log trigger config from the given raw data
// TODO: tests
func (rp *evmRegistryPackerV21) UnpackLogTriggerConfig(raw string) (iregistry21.KeeperRegistryBase21LogTriggerConfig, error) {
	var cfg iregistry21.KeeperRegistryBase21LogTriggerConfig
	b, err := hexutil.Decode(raw)
	if err != nil {
		return cfg, err
	}

	out, err := rp.abi.Methods["getLogTriggerConfig"].Outputs.UnpackValues(b)
	if err != nil {
		return cfg, fmt.Errorf("%w: unpack getLogTriggerConfig return: %s", err, raw)
	}

	converted, ok := abi.ConvertType(out[0], new(iregistry21.KeeperRegistryBase21LogTriggerConfig)).(*iregistry21.KeeperRegistryBase21LogTriggerConfig)
	if !ok {
		return cfg, fmt.Errorf("failed to convert type")
	}
	return *converted, nil
}