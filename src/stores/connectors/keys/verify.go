package keys

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/crypto"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
)

func (c Connector) Verify(_ context.Context, pubKey, data, sig []byte, algo *entities.Algorithm) error {
	logger := c.logger.With("pub_key", pubKey, "curve", algo.EllipticCurve, "signing_algorithm", algo.Type)

	var err error
	var verified bool
	switch {
	case algo.EllipticCurve == entities.Secp256k1 && algo.Type == entities.Ecdsa:
		verified, err = crypto.VerifyECDSASignature(pubKey, data, sig)
	case algo.EllipticCurve == entities.Bn254 && algo.Type == entities.Eddsa:
		verified, err = crypto.VerifyEDDSASignature(pubKey, data, sig)
	default:
		errMessage := "unsupported signing algorithm and elliptic curve combination"
		logger.Error(errMessage)
		return errors.NotSupportedError(errMessage)
	}
	if err != nil {
		errMessage := "failed to verify signature"
		logger.WithError(err).Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	if !verified {
		errMessage := "signature does not belong to the specified public key"
		logger.Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	c.logger.Debug("data verified successfully")
	return nil
}
