package json

import (
	"github.com/consensys/quorum-key-manager/src/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

func isHexAddress(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		return common.IsHexAddress(fl.Field().String())
	}

	return true
}

func isCurve(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		switch fl.Field().String() {
		case string(entities.Secp256k1), string(entities.Babyjubjub):
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
		case string(entities.Ecdsa), string(entities.Eddsa):
			return true
		default:
			return false
		}
	}

	return true
}

func isAliasKind(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		switch fl.Field().String() {
		case entities.AliasKindString, entities.AliasKindArray:
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
	_ = validate.RegisterValidation("isAliasKind", isAliasKind)
}

func getValidator() *validator.Validate {
	return validate
}
