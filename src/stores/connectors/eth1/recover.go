package eth1

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

func (c Connector) ECRecover(_ context.Context, data, sig []byte) (string, error) {
	pubKey, err := crypto.SigToPub(crypto.Keccak256(data), sig)
	if err != nil {
		errMessage := "failed to recover public key, please verify your signature and payload"
		c.logger.WithError(err).Error(errMessage)
		return "", errors.InvalidParameterError(errMessage)
	}

	c.logger.Debug("ethereum account recovered successfully from signature")
	return crypto.PubkeyToAddress(*pubKey).Hex(), nil
}
