package eth1

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (c Connector) ECRecover(_ context.Context, data, sig []byte) (ethcommon.Address, error) {
	pubKey, err := crypto.SigToPub(crypto.Keccak256(data), sig)
	if err != nil {
		errMessage := "failed to recover public key, please verify your signature and payload"
		c.logger.WithError(err).Error(errMessage)
		return ethcommon.Address{}, errors.InvalidParameterError(errMessage)
	}

	c.logger.Debug("ethereum account recovered successfully from signature")
	return crypto.PubkeyToAddress(*pubKey), nil
}
