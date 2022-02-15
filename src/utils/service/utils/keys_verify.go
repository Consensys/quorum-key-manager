package utils

import (
	"github.com/consensys/quorum-key-manager/pkg/crypto/ecdsa"
	"github.com/consensys/quorum-key-manager/pkg/crypto/eddsa"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (u *Utilities) Verify(pubKey, data, sig []byte, algo *entities.Algorithm) error {
	logger := u.logger.With("pub_key", pubKey, "curve", algo.EllipticCurve, "signing_algorithm", algo.Type)

	var err error
	var verified bool
	switch {
	case algo.EllipticCurve == entities.Secp256k1 && algo.Type == entities.Ecdsa:
		verified, err = ecdsa.VerifySecp256k1Signature(pubKey, data, sig)
	case algo.EllipticCurve == entities.Babyjubjub && algo.Type == entities.Eddsa:
		verified, err = eddsa.VerifyBabyJubJubSignature(pubKey, data, sig)
	case algo.EllipticCurve == entities.Curve25519 && algo.Type == entities.Eddsa:
		verified, err = eddsa.VerifyED25519Signature(pubKey, data, sig)
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

	u.logger.Debug("data verified successfully")
	return nil
}
