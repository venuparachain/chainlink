// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	big "math/big"

	assets "github.com/smartcontractkit/chainlink/v2/core/assets"

	common "github.com/ethereum/go-ethereum/common"

	config "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"

	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"

	ethkey "github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"

	mock "github.com/stretchr/testify/mock"

	time "time"

	uuid "github.com/google/uuid"

	zapcore "go.uber.org/zap/zapcore"
)

// ChainScopedConfig is an autogenerated mock type for the ChainScopedConfig type
type ChainScopedConfig struct {
	mock.Mock
}

// AppID provides a mock function with given fields:
func (_m *ChainScopedConfig) AppID() uuid.UUID {
	ret := _m.Called()

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func() uuid.UUID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	return r0
}

// AuditLogger provides a mock function with given fields:
func (_m *ChainScopedConfig) AuditLogger() coreconfig.AuditLogger {
	ret := _m.Called()

	var r0 coreconfig.AuditLogger
	if rf, ok := ret.Get(0).(func() coreconfig.AuditLogger); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.AuditLogger)
		}
	}

	return r0
}

// AutoCreateKey provides a mock function with given fields:
func (_m *ChainScopedConfig) AutoCreateKey() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// AutoPprof provides a mock function with given fields:
func (_m *ChainScopedConfig) AutoPprof() coreconfig.AutoPprof {
	ret := _m.Called()

	var r0 coreconfig.AutoPprof
	if rf, ok := ret.Get(0).(func() coreconfig.AutoPprof); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.AutoPprof)
		}
	}

	return r0
}

// BlockBackfillDepth provides a mock function with given fields:
func (_m *ChainScopedConfig) BlockBackfillDepth() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// BlockBackfillSkip provides a mock function with given fields:
func (_m *ChainScopedConfig) BlockBackfillSkip() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// BlockEmissionIdleWarningThreshold provides a mock function with given fields:
func (_m *ChainScopedConfig) BlockEmissionIdleWarningThreshold() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// ChainID provides a mock function with given fields:
func (_m *ChainScopedConfig) ChainID() *big.Int {
	ret := _m.Called()

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func() *big.Int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	return r0
}

// ChainType provides a mock function with given fields:
func (_m *ChainScopedConfig) ChainType() coreconfig.ChainType {
	ret := _m.Called()

	var r0 coreconfig.ChainType
	if rf, ok := ret.Get(0).(func() coreconfig.ChainType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(coreconfig.ChainType)
	}

	return r0
}

// CosmosEnabled provides a mock function with given fields:
func (_m *ChainScopedConfig) CosmosEnabled() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Database provides a mock function with given fields:
func (_m *ChainScopedConfig) Database() coreconfig.Database {
	ret := _m.Called()

	var r0 coreconfig.Database
	if rf, ok := ret.Get(0).(func() coreconfig.Database); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Database)
		}
	}

	return r0
}

// DefaultChainID provides a mock function with given fields:
func (_m *ChainScopedConfig) DefaultChainID() *big.Int {
	ret := _m.Called()

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func() *big.Int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	return r0
}

// EVM provides a mock function with given fields:
func (_m *ChainScopedConfig) EVM() config.EVM {
	ret := _m.Called()

	var r0 config.EVM
	if rf, ok := ret.Get(0).(func() config.EVM); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(config.EVM)
		}
	}

	return r0
}

// EVMEnabled provides a mock function with given fields:
func (_m *ChainScopedConfig) EVMEnabled() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// EVMRPCEnabled provides a mock function with given fields:
func (_m *ChainScopedConfig) EVMRPCEnabled() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// EvmEIP1559DynamicFees provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmEIP1559DynamicFees() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// EvmFinalityDepth provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmFinalityDepth() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// EvmGasBumpPercent provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasBumpPercent() uint16 {
	ret := _m.Called()

	var r0 uint16
	if rf, ok := ret.Get(0).(func() uint16); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint16)
	}

	return r0
}

// EvmGasBumpThreshold provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasBumpThreshold() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// EvmGasBumpTxDepth provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasBumpTxDepth() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// EvmGasBumpWei provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasBumpWei() *assets.Wei {
	ret := _m.Called()

	var r0 *assets.Wei
	if rf, ok := ret.Get(0).(func() *assets.Wei); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Wei)
		}
	}

	return r0
}

// EvmGasFeeCapDefault provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasFeeCapDefault() *assets.Wei {
	ret := _m.Called()

	var r0 *assets.Wei
	if rf, ok := ret.Get(0).(func() *assets.Wei); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Wei)
		}
	}

	return r0
}

// EvmGasLimitDRJobType provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasLimitDRJobType() *uint32 {
	ret := _m.Called()

	var r0 *uint32
	if rf, ok := ret.Get(0).(func() *uint32); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uint32)
		}
	}

	return r0
}

// EvmGasLimitDefault provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasLimitDefault() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// EvmGasLimitFMJobType provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasLimitFMJobType() *uint32 {
	ret := _m.Called()

	var r0 *uint32
	if rf, ok := ret.Get(0).(func() *uint32); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uint32)
		}
	}

	return r0
}

// EvmGasLimitKeeperJobType provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasLimitKeeperJobType() *uint32 {
	ret := _m.Called()

	var r0 *uint32
	if rf, ok := ret.Get(0).(func() *uint32); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uint32)
		}
	}

	return r0
}

// EvmGasLimitMax provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasLimitMax() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// EvmGasLimitMultiplier provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasLimitMultiplier() float32 {
	ret := _m.Called()

	var r0 float32
	if rf, ok := ret.Get(0).(func() float32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(float32)
	}

	return r0
}

// EvmGasLimitOCR2JobType provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasLimitOCR2JobType() *uint32 {
	ret := _m.Called()

	var r0 *uint32
	if rf, ok := ret.Get(0).(func() *uint32); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uint32)
		}
	}

	return r0
}

// EvmGasLimitOCRJobType provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasLimitOCRJobType() *uint32 {
	ret := _m.Called()

	var r0 *uint32
	if rf, ok := ret.Get(0).(func() *uint32); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uint32)
		}
	}

	return r0
}

// EvmGasLimitTransfer provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasLimitTransfer() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// EvmGasLimitVRFJobType provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasLimitVRFJobType() *uint32 {
	ret := _m.Called()

	var r0 *uint32
	if rf, ok := ret.Get(0).(func() *uint32); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uint32)
		}
	}

	return r0
}

// EvmGasPriceDefault provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasPriceDefault() *assets.Wei {
	ret := _m.Called()

	var r0 *assets.Wei
	if rf, ok := ret.Get(0).(func() *assets.Wei); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Wei)
		}
	}

	return r0
}

// EvmGasTipCapDefault provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasTipCapDefault() *assets.Wei {
	ret := _m.Called()

	var r0 *assets.Wei
	if rf, ok := ret.Get(0).(func() *assets.Wei); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Wei)
		}
	}

	return r0
}

// EvmGasTipCapMinimum provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmGasTipCapMinimum() *assets.Wei {
	ret := _m.Called()

	var r0 *assets.Wei
	if rf, ok := ret.Get(0).(func() *assets.Wei); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Wei)
		}
	}

	return r0
}

// EvmHeadTrackerHistoryDepth provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmHeadTrackerHistoryDepth() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// EvmHeadTrackerMaxBufferSize provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmHeadTrackerMaxBufferSize() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// EvmHeadTrackerSamplingInterval provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmHeadTrackerSamplingInterval() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// EvmLogBackfillBatchSize provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmLogBackfillBatchSize() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// EvmLogKeepBlocksDepth provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmLogKeepBlocksDepth() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// EvmLogPollInterval provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmLogPollInterval() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// EvmMaxGasPriceWei provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmMaxGasPriceWei() *assets.Wei {
	ret := _m.Called()

	var r0 *assets.Wei
	if rf, ok := ret.Get(0).(func() *assets.Wei); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Wei)
		}
	}

	return r0
}

// EvmMinGasPriceWei provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmMinGasPriceWei() *assets.Wei {
	ret := _m.Called()

	var r0 *assets.Wei
	if rf, ok := ret.Get(0).(func() *assets.Wei); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Wei)
		}
	}

	return r0
}

// EvmNonceAutoSync provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmNonceAutoSync() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// EvmRPCDefaultBatchSize provides a mock function with given fields:
func (_m *ChainScopedConfig) EvmRPCDefaultBatchSize() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// Explorer provides a mock function with given fields:
func (_m *ChainScopedConfig) Explorer() coreconfig.Explorer {
	ret := _m.Called()

	var r0 coreconfig.Explorer
	if rf, ok := ret.Get(0).(func() coreconfig.Explorer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Explorer)
		}
	}

	return r0
}

// Feature provides a mock function with given fields:
func (_m *ChainScopedConfig) Feature() coreconfig.Feature {
	ret := _m.Called()

	var r0 coreconfig.Feature
	if rf, ok := ret.Get(0).(func() coreconfig.Feature); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Feature)
		}
	}

	return r0
}

// FlagsContractAddress provides a mock function with given fields:
func (_m *ChainScopedConfig) FlagsContractAddress() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// FluxMonitor provides a mock function with given fields:
func (_m *ChainScopedConfig) FluxMonitor() coreconfig.FluxMonitor {
	ret := _m.Called()

	var r0 coreconfig.FluxMonitor
	if rf, ok := ret.Get(0).(func() coreconfig.FluxMonitor); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.FluxMonitor)
		}
	}

	return r0
}

// GasEstimatorMode provides a mock function with given fields:
func (_m *ChainScopedConfig) GasEstimatorMode() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Insecure provides a mock function with given fields:
func (_m *ChainScopedConfig) Insecure() coreconfig.Insecure {
	ret := _m.Called()

	var r0 coreconfig.Insecure
	if rf, ok := ret.Get(0).(func() coreconfig.Insecure); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Insecure)
		}
	}

	return r0
}

// InsecureFastScrypt provides a mock function with given fields:
func (_m *ChainScopedConfig) InsecureFastScrypt() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// JobPipeline provides a mock function with given fields:
func (_m *ChainScopedConfig) JobPipeline() coreconfig.JobPipeline {
	ret := _m.Called()

	var r0 coreconfig.JobPipeline
	if rf, ok := ret.Get(0).(func() coreconfig.JobPipeline); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.JobPipeline)
		}
	}

	return r0
}

// Keeper provides a mock function with given fields:
func (_m *ChainScopedConfig) Keeper() coreconfig.Keeper {
	ret := _m.Called()

	var r0 coreconfig.Keeper
	if rf, ok := ret.Get(0).(func() coreconfig.Keeper); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Keeper)
		}
	}

	return r0
}

// KeySpecificMaxGasPriceWei provides a mock function with given fields: addr
func (_m *ChainScopedConfig) KeySpecificMaxGasPriceWei(addr common.Address) *assets.Wei {
	ret := _m.Called(addr)

	var r0 *assets.Wei
	if rf, ok := ret.Get(0).(func(common.Address) *assets.Wei); ok {
		r0 = rf(addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Wei)
		}
	}

	return r0
}

// LinkContractAddress provides a mock function with given fields:
func (_m *ChainScopedConfig) LinkContractAddress() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Log provides a mock function with given fields:
func (_m *ChainScopedConfig) Log() coreconfig.Log {
	ret := _m.Called()

	var r0 coreconfig.Log
	if rf, ok := ret.Get(0).(func() coreconfig.Log); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Log)
		}
	}

	return r0
}

// LogConfiguration provides a mock function with given fields: log
func (_m *ChainScopedConfig) LogConfiguration(log coreconfig.LogfFn) {
	_m.Called(log)
}

// Mercury provides a mock function with given fields:
func (_m *ChainScopedConfig) Mercury() coreconfig.Mercury {
	ret := _m.Called()

	var r0 coreconfig.Mercury
	if rf, ok := ret.Get(0).(func() coreconfig.Mercury); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Mercury)
		}
	}

	return r0
}

// MinIncomingConfirmations provides a mock function with given fields:
func (_m *ChainScopedConfig) MinIncomingConfirmations() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// MinimumContractPayment provides a mock function with given fields:
func (_m *ChainScopedConfig) MinimumContractPayment() *assets.Link {
	ret := _m.Called()

	var r0 *assets.Link
	if rf, ok := ret.Get(0).(func() *assets.Link); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Link)
		}
	}

	return r0
}

// NodeNoNewHeadsThreshold provides a mock function with given fields:
func (_m *ChainScopedConfig) NodeNoNewHeadsThreshold() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// NodePollFailureThreshold provides a mock function with given fields:
func (_m *ChainScopedConfig) NodePollFailureThreshold() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// NodePollInterval provides a mock function with given fields:
func (_m *ChainScopedConfig) NodePollInterval() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// NodeSelectionMode provides a mock function with given fields:
func (_m *ChainScopedConfig) NodeSelectionMode() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NodeSyncThreshold provides a mock function with given fields:
func (_m *ChainScopedConfig) NodeSyncThreshold() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// OCR2AutomationGasLimit provides a mock function with given fields:
func (_m *ChainScopedConfig) OCR2AutomationGasLimit() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// OCR2BlockchainTimeout provides a mock function with given fields:
func (_m *ChainScopedConfig) OCR2BlockchainTimeout() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// OCR2CaptureEATelemetry provides a mock function with given fields:
func (_m *ChainScopedConfig) OCR2CaptureEATelemetry() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// OCR2ContractConfirmations provides a mock function with given fields:
func (_m *ChainScopedConfig) OCR2ContractConfirmations() uint16 {
	ret := _m.Called()

	var r0 uint16
	if rf, ok := ret.Get(0).(func() uint16); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint16)
	}

	return r0
}

// OCR2ContractPollInterval provides a mock function with given fields:
func (_m *ChainScopedConfig) OCR2ContractPollInterval() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// OCR2ContractSubscribeInterval provides a mock function with given fields:
func (_m *ChainScopedConfig) OCR2ContractSubscribeInterval() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// OCR2ContractTransmitterTransmitTimeout provides a mock function with given fields:
func (_m *ChainScopedConfig) OCR2ContractTransmitterTransmitTimeout() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// OCR2DatabaseTimeout provides a mock function with given fields:
func (_m *ChainScopedConfig) OCR2DatabaseTimeout() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// OCR2DefaultTransactionQueueDepth provides a mock function with given fields:
func (_m *ChainScopedConfig) OCR2DefaultTransactionQueueDepth() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// OCR2Enabled provides a mock function with given fields:
func (_m *ChainScopedConfig) OCR2Enabled() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// OCR2KeyBundleID provides a mock function with given fields:
func (_m *ChainScopedConfig) OCR2KeyBundleID() (string, error) {
	ret := _m.Called()

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OCR2SimulateTransactions provides a mock function with given fields:
func (_m *ChainScopedConfig) OCR2SimulateTransactions() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// OCR2TraceLogging provides a mock function with given fields:
func (_m *ChainScopedConfig) OCR2TraceLogging() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// OCRBlockchainTimeout provides a mock function with given fields:
func (_m *ChainScopedConfig) OCRBlockchainTimeout() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// OCRCaptureEATelemetry provides a mock function with given fields:
func (_m *ChainScopedConfig) OCRCaptureEATelemetry() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// OCRContractConfirmations provides a mock function with given fields:
func (_m *ChainScopedConfig) OCRContractConfirmations() uint16 {
	ret := _m.Called()

	var r0 uint16
	if rf, ok := ret.Get(0).(func() uint16); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint16)
	}

	return r0
}

// OCRContractPollInterval provides a mock function with given fields:
func (_m *ChainScopedConfig) OCRContractPollInterval() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// OCRContractSubscribeInterval provides a mock function with given fields:
func (_m *ChainScopedConfig) OCRContractSubscribeInterval() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// OCRContractTransmitterTransmitTimeout provides a mock function with given fields:
func (_m *ChainScopedConfig) OCRContractTransmitterTransmitTimeout() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// OCRDatabaseTimeout provides a mock function with given fields:
func (_m *ChainScopedConfig) OCRDatabaseTimeout() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// OCRDefaultTransactionQueueDepth provides a mock function with given fields:
func (_m *ChainScopedConfig) OCRDefaultTransactionQueueDepth() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// OCREnabled provides a mock function with given fields:
func (_m *ChainScopedConfig) OCREnabled() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// OCRKeyBundleID provides a mock function with given fields:
func (_m *ChainScopedConfig) OCRKeyBundleID() (string, error) {
	ret := _m.Called()

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OCRObservationGracePeriod provides a mock function with given fields:
func (_m *ChainScopedConfig) OCRObservationGracePeriod() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// OCRObservationTimeout provides a mock function with given fields:
func (_m *ChainScopedConfig) OCRObservationTimeout() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// OCRSimulateTransactions provides a mock function with given fields:
func (_m *ChainScopedConfig) OCRSimulateTransactions() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// OCRTraceLogging provides a mock function with given fields:
func (_m *ChainScopedConfig) OCRTraceLogging() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// OCRTransmitterAddress provides a mock function with given fields:
func (_m *ChainScopedConfig) OCRTransmitterAddress() (ethkey.EIP55Address, error) {
	ret := _m.Called()

	var r0 ethkey.EIP55Address
	var r1 error
	if rf, ok := ret.Get(0).(func() (ethkey.EIP55Address, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() ethkey.EIP55Address); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(ethkey.EIP55Address)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OperatorFactoryAddress provides a mock function with given fields:
func (_m *ChainScopedConfig) OperatorFactoryAddress() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// P2P provides a mock function with given fields:
func (_m *ChainScopedConfig) P2P() coreconfig.P2P {
	ret := _m.Called()

	var r0 coreconfig.P2P
	if rf, ok := ret.Get(0).(func() coreconfig.P2P); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.P2P)
		}
	}

	return r0
}

// Password provides a mock function with given fields:
func (_m *ChainScopedConfig) Password() coreconfig.Password {
	ret := _m.Called()

	var r0 coreconfig.Password
	if rf, ok := ret.Get(0).(func() coreconfig.Password); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Password)
		}
	}

	return r0
}

// Prometheus provides a mock function with given fields:
func (_m *ChainScopedConfig) Prometheus() coreconfig.Prometheus {
	ret := _m.Called()

	var r0 coreconfig.Prometheus
	if rf, ok := ret.Get(0).(func() coreconfig.Prometheus); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Prometheus)
		}
	}

	return r0
}

// Pyroscope provides a mock function with given fields:
func (_m *ChainScopedConfig) Pyroscope() coreconfig.Pyroscope {
	ret := _m.Called()

	var r0 coreconfig.Pyroscope
	if rf, ok := ret.Get(0).(func() coreconfig.Pyroscope); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Pyroscope)
		}
	}

	return r0
}

// RootDir provides a mock function with given fields:
func (_m *ChainScopedConfig) RootDir() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Sentry provides a mock function with given fields:
func (_m *ChainScopedConfig) Sentry() coreconfig.Sentry {
	ret := _m.Called()

	var r0 coreconfig.Sentry
	if rf, ok := ret.Get(0).(func() coreconfig.Sentry); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Sentry)
		}
	}

	return r0
}

// SetLogLevel provides a mock function with given fields: lvl
func (_m *ChainScopedConfig) SetLogLevel(lvl zapcore.Level) error {
	ret := _m.Called(lvl)

	var r0 error
	if rf, ok := ret.Get(0).(func(zapcore.Level) error); ok {
		r0 = rf(lvl)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetLogSQL provides a mock function with given fields: logSQL
func (_m *ChainScopedConfig) SetLogSQL(logSQL bool) {
	_m.Called(logSQL)
}

// SetPasswords provides a mock function with given fields: keystore, vrf
func (_m *ChainScopedConfig) SetPasswords(keystore *string, vrf *string) {
	_m.Called(keystore, vrf)
}

// ShutdownGracePeriod provides a mock function with given fields:
func (_m *ChainScopedConfig) ShutdownGracePeriod() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// SolanaEnabled provides a mock function with given fields:
func (_m *ChainScopedConfig) SolanaEnabled() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// StarkNetEnabled provides a mock function with given fields:
func (_m *ChainScopedConfig) StarkNetEnabled() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// TelemetryIngress provides a mock function with given fields:
func (_m *ChainScopedConfig) TelemetryIngress() coreconfig.TelemetryIngress {
	ret := _m.Called()

	var r0 coreconfig.TelemetryIngress
	if rf, ok := ret.Get(0).(func() coreconfig.TelemetryIngress); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.TelemetryIngress)
		}
	}

	return r0
}

// Threshold provides a mock function with given fields:
func (_m *ChainScopedConfig) Threshold() coreconfig.Threshold {
	ret := _m.Called()

	var r0 coreconfig.Threshold
	if rf, ok := ret.Get(0).(func() coreconfig.Threshold); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Threshold)
		}
	}

	return r0
}

// Validate provides a mock function with given fields:
func (_m *ChainScopedConfig) Validate() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateDB provides a mock function with given fields:
func (_m *ChainScopedConfig) ValidateDB() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WebServer provides a mock function with given fields:
func (_m *ChainScopedConfig) WebServer() coreconfig.WebServer {
	ret := _m.Called()

	var r0 coreconfig.WebServer
	if rf, ok := ret.Get(0).(func() coreconfig.WebServer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.WebServer)
		}
	}

	return r0
}

type mockConstructorTestingTNewChainScopedConfig interface {
	mock.TestingT
	Cleanup(func())
}

// NewChainScopedConfig creates a new instance of ChainScopedConfig. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewChainScopedConfig(t mockConstructorTestingTNewChainScopedConfig) *ChainScopedConfig {
	mock := &ChainScopedConfig{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
