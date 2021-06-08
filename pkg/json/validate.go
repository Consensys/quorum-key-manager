package json

import (
	entities2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
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
		case string(entities2.Secp256k1), string(entities2.Bn254):
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

func init() {
	if validate != nil {
		return
	}

	validate = validator.New()
	_ = validate.RegisterValidation("isHexAddress", isHexAddress)
	_ = validate.RegisterValidation("isCurve", isCurve)
	_ = validate.RegisterValidation("isSigningAlgorithm", isSigningAlgorithm)
}

func getValidator() *validator.Validate {
	return validate
}
