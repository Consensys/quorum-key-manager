package utils

import (
	"github.com/consensys/quorum-key-manager/pkg/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (c Connector) ECRecover(data, sig []byte) (ethcommon.Address, error) {
	pubKey, err := crypto.SigToPub(crypto.Keccak256(data), sig)
	if err != nil {
		c.logger.WithError(err).Error("failed to recover public key, please verify your signature and payload")
		return ethcommon.Address{}, errors.InvalidParameterError(err.Error())
	}

	c.logger.Debug("ethereum account recovered successfully")
	return crypto.PubkeyToAddress(*pubKey), nil
}
