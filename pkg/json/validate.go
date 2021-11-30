package json

import (
	entities2 "github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

func isHexAddress(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		return ethcommon.IsHexAddress(fl.Field().String())
	}

	return true
}

func isCurve(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		switch fl.Field().String() {
		case string(entities2.Secp256k1), string(entities2.Babyjubjub):
			return true
		default:
			return false
		}
	}

	return true
}

func isSigningAlgorithm(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		switch fl.Field().String() {
		case string(entities2.Ecdsa), string(entities2.Eddsa):
			return true
		default:
			return false
		}
	}

	return true
}

func isManifestKind(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		switch fl.Field().String() {
		case entities2.StoreKind, entities2.NodeKind, entities2.RoleKind, entities2.VaultKind:
			return true
		default:
			return false
		}
	}

	return true
}

func isStoreType(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		switch fl.Field().String() {
		case entities.SecretStoreType, entities.KeyStoreType, entities.EthereumStoreType:
			return true
		default:
			return false
		}
	}

	return true
}

func init() {
	if validate != nil {
		return
	}

	validate = validator.New()
	_ = validate.RegisterValidation("isHexAddress", isHexAddress)
	_ = validate.RegisterValidation("isCurve", isCurve)
	_ = validate.RegisterValidation("isSigningAlgorithm", isSigningAlgorithm)
	_ = validate.RegisterValidation("isManifestKind", isManifestKind)
	_ = validate.RegisterValidation("isStoreType", isStoreType)
}

func getValidator() *validator.Validate {
	return validate
}
