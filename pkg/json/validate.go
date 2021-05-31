package json

import (
	"encoding/base64"
	"math/big"
	"reflect"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-playground/validator/v10"
)

var (
	validate      *validator.Validate
	StringPtrType = reflect.TypeOf(new(string))
	StringType    = reflect.TypeOf("")
)

func isHex(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		_, err := hexutil.Decode(fl.Field().String())
		return err == nil
	}

	return true
}

func isBase64(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		_, err := base64.URLEncoding.DecodeString(fl.Field().String())
		return err == nil
	}

	return true
}

func isHexAddress(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		return ethcommon.IsHexAddress(fl.Field().String())
	}

	return true
}

func isBig(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		_, ok := new(big.Int).SetString(fl.Field().String(), 10)
		return ok
	}

	return true
}

func isDuration(fl validator.FieldLevel) bool {
	_, err := convDuration(fl)
	return err == nil
}

func convDuration(fl validator.FieldLevel) (time.Duration, error) {
	switch fl.Field().Type() {
	case StringPtrType:
		val := fl.Field().Interface().(*string)
		if val != nil {
			return time.ParseDuration(*val)
		}
		return time.Duration(0), nil
	case StringType:
		if fl.Field().String() != "" {
			return time.ParseDuration(fl.Field().String())
		}
		return time.Duration(0), nil
	default:
		return time.Duration(0), nil
	}
}

func isCurve(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		switch fl.Field().String() {
		case string(entities.Secp256k1), string(entities.Bn254):
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

func init() {
	if validate != nil {
		return
	}

	validate = validator.New()
	_ = validate.RegisterValidation("isHex", isHex)
	_ = validate.RegisterValidation("isBig", isBig)
	_ = validate.RegisterValidation("isHexAddress", isHexAddress)
	_ = validate.RegisterValidation("isDuration", isDuration)
	_ = validate.RegisterValidation("isCurve", isCurve)
	_ = validate.RegisterValidation("isSigningAlgorithm", isSigningAlgorithm)
	_ = validate.RegisterValidation("isBase64", isBase64)
}

func getValidator() *validator.Validate {
	return validate
}
