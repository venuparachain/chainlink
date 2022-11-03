package coordinator

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/recovery_beacon"
)

type topics struct {
	configSetTopic common.Hash
}

func newTopics() topics {
	return topics{
		configSetTopic: recovery_beacon.RecoveryBeaconConfigSet{}.Topic(),
	}
}
