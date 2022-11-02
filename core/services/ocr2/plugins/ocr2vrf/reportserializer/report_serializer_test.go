package reportserializer_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	"github.com/smartcontractkit/ocr2vrf/types"

	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/reportserializer"
)

func Test_Serialize_Deserialize(t *testing.T) {
	altbn128Suite := &altbn_128.PairingSuite{}
	reportSerializer := reportserializer.NewReportSerializer(altbn128Suite.G1())

	unserializedReport := types.AbstractReport{
		Recoverer:        common.Address{},
		AccountToRecover: common.Address{},
	}
	r, err := reportSerializer.SerializeReport(unserializedReport)
	require.NoError(t, err)
	require.Equal(t, uint(len(r)), reportSerializer.ReportLength(unserializedReport))
	// TODO: Add deserialization after this point to verify.
}
