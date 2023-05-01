package evm

import (
	"errors"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"
)

// TODO: move

type UpkeepConfig interface {
	// Validate performs static validation of the fields
	Validate() error
	Decode([]byte) error
	Encode() ([]byte, error)
}

func decodeUpkeepStateEvent(cfg UpkeepConfig, abilog generated.AbigenLog) (types.UpkeepIdentifier, error) {
	var uid types.UpkeepIdentifier
	var err error
	switch l := abilog.(type) {
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepCanceled:
		uid = types.UpkeepIdentifier(l.Id.String())
		err = cfg.Decode(l.Raw.Data)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepPaused:
		uid = types.UpkeepIdentifier(l.Id.String())
		err = cfg.Decode(l.Raw.Data)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepRegistered:
		uid = types.UpkeepIdentifier(l.Id.String())
		err = cfg.Decode(l.Raw.Data)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepReceived:
		uid = types.UpkeepIdentifier(l.Id.String())
		err = cfg.Decode(l.Raw.Data)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepUnpaused:
		uid = types.UpkeepIdentifier(l.Id.String())
		err = cfg.Decode(l.Raw.Data)
	default:
		return nil, errors.New("unsupported type")
	}
	return uid, err
}
