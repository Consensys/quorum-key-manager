package utils

import (
	"github.com/consensys/quorum-key-manager/pkg/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (u *Utilities) ECRecover(data, sig []byte) (ethcommon.Address, error) {
	pubKey, err := crypto.SigToPub(crypto.Keccak256(data), sig)
	if err != nil {
		u.logger.WithError(err).Error("failed to recover public key, please verify your signature and payload")
		return ethcommon.Address{}, errors.InvalidParameterError(err.Error())
	}

	u.logger.Debug("ethereum account recovered successfully")
	return crypto.PubkeyToAddress(*pubKey), nil
}
