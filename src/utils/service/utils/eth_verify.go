package utils

import (
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/ethereum"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core"
)

func (u *Utilities) VerifyMessage(addr ethcommon.Address, data, sig []byte) error {
	err := u.verifyHomestead(addr, ethereum.GetEIP191EncodedData(data), sig)
	if err != nil {
		return err
	}

	u.logger.Debug("message signature verified successfully")
	return nil
}

func (u *Utilities) VerifyTypedData(addr ethcommon.Address, typedData *core.TypedData, sig []byte) error {
	encodedData, err := ethereum.GetEIP712EncodedData(typedData)
	if err != nil {
		errMessage := "failed to generate EIP-712 encoded data"
		u.logger.WithError(err).Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	err = u.verifyHomestead(addr, encodedData, sig)
	if err != nil {
		return err
	}

	u.logger.Debug("typed data signature verified successfully")
	return nil
}

func (u *Utilities) verifyHomestead(addr ethcommon.Address, data, sig []byte) error {
	sigLength := len(sig)
	if sigLength != crypto.SignatureLength {
		errMessage := "signature must be exactly 65 bytes"
		u.logger.With("signature_length", sigLength).Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	sig[crypto.RecoveryIDOffset] -= 27

	recoveredAddress, err := u.ECRecover(data, sig)
	if err != nil {
		return err
	}

	if addr.Hex() != recoveredAddress.Hex() {
		errMessage := "failed to verify signature: recovered address does not match the expected one or payload is malformed"
		u.logger.WithError(err).With("address", addr.Hex(), "recovered_address", recoveredAddress.Hex()).Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	return nil
}
