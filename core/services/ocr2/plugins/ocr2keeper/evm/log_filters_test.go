package evm

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/stretchr/testify/require"
)

func TestLogFiltersProvider_Register(t *testing.T) {
	tests := []struct {
		name         string
		errored      bool
		upkeepID     *big.Int
		upkeepCfg    *LogTriggerUpkeepConfig
		rawUpkeepCfg []byte
	}{
		{
			"happy flow",
			false,
			big.NewInt(111),
			&LogTriggerUpkeepConfig{
				Address: "0x" + string(common.LeftPadBytes([]byte{}, 20)),
				Topic:   string(common.LeftPadBytes([]byte{}, 32)),
			},
			nil,
		},
		{
			"invalid encoding",
			true,
			big.NewInt(111),
			nil,
			[]byte("1111"),
		},
		{
			"invalid config",
			true,
			big.NewInt(111),
			&LogTriggerUpkeepConfig{
				// missing 0x prefix
				Address: string(common.LeftPadBytes([]byte{}, 20)),
				Topic:   string(common.LeftPadBytes([]byte{}, 32)),
			},
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mp := new(mocks.LogPoller)
			if tc.upkeepCfg != nil {
				mp.On("RegisterFilter", newLogFilter(tc.upkeepCfg, tc.upkeepID.String())).Return(nil)
			}
			lfp := newLogFiltersProvider(mp)
			if len(tc.rawUpkeepCfg) == 0 {
				require.NotNil(t, tc.upkeepCfg)
				raw, err := tc.upkeepCfg.Encode()
				require.NoError(t, err)
				tc.rawUpkeepCfg = raw
			}
			err := lfp.Register(tc.upkeepID.String(), tc.rawUpkeepCfg)
			if tc.errored {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
