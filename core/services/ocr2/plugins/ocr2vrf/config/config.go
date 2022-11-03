package config

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	dkgconfig "github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/dkg/config"
)

// PluginConfig contains custom arguments for the OCR2Recovery plugin.
//
// The OCR2Recovery plugin runs a DKG under the hood, so it will need both
// DKG and OCR2Recovery configuration fields.
//
// The DKG contract address is provided in the plugin configuration,
// however the OCR2Recovery contract address is provided in the OCR2 job spec
// under the 'contractID' key.
type PluginConfig struct {
	// DKG configuration fields.
	DKGEncryptionPublicKey string `json:"dkgEncryptionPublicKey"`
	DKGSigningPublicKey    string `json:"dkgSigningPublicKey"`
	DKGKeyID               string `json:"dkgKeyID"`
	DKGContractAddress     string `json:"dkgContractAddress"`
}

// ValidatePluginConfig validates that the given OCR2Recovery plugin configuration is correct.
func ValidatePluginConfig(config PluginConfig, dkgSignKs keystore.DKGSign, dkgEncryptKs keystore.DKGEncrypt) error {
	err := dkgconfig.ValidatePluginConfig(dkgconfig.PluginConfig{
		EncryptionPublicKey: config.DKGEncryptionPublicKey,
		SigningPublicKey:    config.DKGSigningPublicKey,
		KeyID:               config.DKGKeyID,
	}, dkgSignKs, dkgEncryptKs)
	if err != nil {
		return err
	}

	// NOTE: a better validation would be to call a method on the on-chain contract pointed to by this
	// address.
	if config.DKGContractAddress == "" {
		return errors.New("dkgContractAddress field must be provided")
	}
	return nil
}
