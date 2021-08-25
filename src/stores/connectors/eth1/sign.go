package eth1

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"math/big"

	authtypes "github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/ethereum"
	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/signer/core"
	"golang.org/x/crypto/sha3"
)

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

func (c Connector) Sign(ctx context.Context, addr common.Address, data []byte) ([]byte, error) {
	logger := c.logger.With("address", addr.Hex())

	signature, err := c.sign(ctx, addr, crypto.Keccak256(data))
	if err != nil {
		return nil, err
	}

	logger.Debug("signed payload successfully")
	return signature, nil
}

func (c Connector) SignMessage(ctx context.Context, addr common.Address, data []byte) ([]byte, error) {
	logger := c.logger.With("address", addr)

	signature, err := c.sign(ctx, addr, crypto.Keccak256([]byte(getEIP191EncodedData(data))))
	if err != nil {
		return nil, err
	}

	logger.Debug("message signed successfully (eip-191)")
	return signature, nil
}

func (c Connector) SignTypedData(ctx context.Context, addr common.Address, typedData *core.TypedData) ([]byte, error) {
	logger := c.logger.With("address", addr.Hex())

	encodedData, err := getEIP712EncodedData(typedData)
	if err != nil {
		errMessage := "failed to format typed data"
		logger.WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	signature, err := c.sign(ctx, addr, crypto.Keccak256([]byte(encodedData)))
	if err != nil {
		return nil, err
	}

	logger.Debug("typed data signed successfully (eip-712) ")
	return signature, nil
}

func (c Connector) SignTransaction(ctx context.Context, addr common.Address, chainID *big.Int, tx *types.Transaction) ([]byte, error) {
	logger := c.logger.With("address", addr.Hex())

	signer := types.NewEIP155Signer(chainID)
	txData := signer.Hash(tx).Bytes()

	signature, err := c.sign(ctx, addr, txData)
	if err != nil {
		return nil, err
	}

	signedTx, err := tx.WithSignature(signer, signature)
	if err != nil {
		errMessage := "failed to set transaction signature"
		c.logger.WithError(err).Error(errMessage)
		return nil, errors.DependencyFailureError(errMessage)
	}

	signedRaw, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		errMessage := "failed to RLP encode signed transaction"
		c.logger.WithError(err).Error(errMessage)
		return nil, errors.EncodingError(errMessage)
	}

	logger.Debug("transaction signed successfully")
	return signedRaw, nil
}

func (c Connector) SignEEA(ctx context.Context, addr common.Address, chainID *big.Int, tx *types.Transaction, args *ethereum.PrivateArgs) ([]byte, error) {
	logger := c.logger.With("address", addr.Hex())

	privateFromEncoded, err := base64.StdEncoding.DecodeString(*args.PrivateFrom)
	if err != nil {
		errMessage := "invalid privateFrom param"
		c.logger.WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	privateRecipientEncoded, err := getEncodedPrivateRecipient(args.PrivacyGroupID, args.PrivateFor)
	if err != nil {
		errMessage := "invalid privacyGroupID or privateFor"
		c.logger.WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	hash, err := eeaHash([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		chainID,
		uint(0),
		uint(0),
		privateFromEncoded,
		privateRecipientEncoded,
		*args.PrivateType,
	})
	if err != nil {
		errMessage := "failed to hash EEA transaction"
		c.logger.WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	signature, err := c.sign(ctx, addr, hash[:])
	if err != nil {
		return nil, err
	}

	signedTx, err := tx.WithSignature(types.NewEIP155Signer(chainID), signature)
	if err != nil {
		errMessage := "failed to set eea transaction signature"
		c.logger.WithError(err).Error(errMessage)
		return nil, errors.DependencyFailureError(errMessage)
	}
	V, R, S := signedTx.RawSignatureValues()

	signedRaw, err := rlp.EncodeToBytes([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		V,
		R,
		S,
		privateFromEncoded,
		privateRecipientEncoded,
		*args.PrivateType,
	})
	if err != nil {
		errMessage := "failed to RLP encode signed eea transaction"
		c.logger.WithError(err).Error(errMessage)
		return nil, errors.EncodingError(errMessage)
	}

	logger.Debug("EEA transaction signed successfully")
	return signedRaw, nil
}

func (c Connector) SignPrivate(ctx context.Context, addr common.Address, tx *quorumtypes.Transaction) ([]byte, error) {
	logger := c.logger.With("address", addr.Hex())

	signer := quorumtypes.QuorumPrivateTxSigner{}
	txData := signer.Hash(tx).Bytes()
	signature, err := c.sign(ctx, addr, txData)
	if err != nil {
		return nil, err
	}

	signedTx, err := tx.WithSignature(signer, signature)
	if err != nil {
		errMessage := "failed to set quorum private transaction signature"
		c.logger.WithError(err).Error(errMessage)
		return nil, errors.DependencyFailureError(errMessage)
	}

	signedRaw, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		errMessage := "failed to RLP encode signed quorum private transaction"
		c.logger.WithError(err).Error(errMessage)
		return nil, errors.EncodingError(errMessage)
	}

	logger.Debug("private transaction signed successfully")
	return signedRaw, nil
}

func (c Connector) sign(ctx context.Context, addr common.Address, data []byte) ([]byte, error) {
	err := c.authorizator.Check(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEth1Account})
	if err != nil {
		return nil, err
	}

	acc, err := c.db.Get(ctx, addr.Hex())
	if err != nil {
		return nil, err
	}

	signature, err := c.store.Sign(ctx, acc.KeyID, data, eth1Algo)
	if err != nil {
		return nil, err
	}

	// Recover the recID, please read: http://coders-errand.com/ecrecover-signature-verification-ethereum/
	for _, recID := range []byte{0, 1} {
		appendedSignature := append(malleabilityECDSASignature(signature), recID)
		var recoveredPubKey *ecdsa.PublicKey
		recoveredPubKey, err = crypto.SigToPub(data, appendedSignature)
		if err != nil {
			errMessage := "failed to recover public key candidate with appended recID"
			c.logger.WithError(err).Error(errMessage, "recID", recID)
			return nil, errors.InvalidParameterError(errMessage)
		}

		if bytes.Equal(crypto.FromECDSAPub(recoveredPubKey), acc.PublicKey) {
			return appendedSignature, nil
		}
	}

	errMessage := "failed to recover public key candidate"
	c.logger.Error(errMessage)
	return nil, errors.InvalidParameterError(errMessage)
}

func eeaHash(object interface{}) (hash common.Hash, err error) {
	hashAlgo := sha3.NewLegacyKeccak256()
	err = rlp.Encode(hashAlgo, object)
	if err != nil {
		return common.Hash{}, err
	}

	hashAlgo.Sum(hash[:0])

	return hash, nil
}

// TODO: Delete usage of unnecessary pointers: https://app.zenhub.com/workspaces/orchestrate-5ea70772b186e10067f57842/issues/consensys/quorum-key-manager/96
func getEncodedPrivateRecipient(privacyGroupID *string, privateFor *[]string) (interface{}, error) {
	var privateRecipientEncoded interface{}
	var err error
	if privacyGroupID != nil {
		privateRecipientEncoded, err = base64.StdEncoding.DecodeString(*privacyGroupID)
		if err != nil {
			return nil, err
		}
	} else {
		var privateForByteSlice [][]byte
		for _, v := range *privateFor {
			b, der := base64.StdEncoding.DecodeString(v)
			if der != nil {
				return nil, err
			}
			privateForByteSlice = append(privateForByteSlice, b)
		}
		privateRecipientEncoded = privateForByteSlice
	}

	return privateRecipientEncoded, nil
}

// Azure generates ECDSA signature that does not prevent malleability
// A malleable signature can be transformed into a new and valid one for a different message or key.
// https://docs.microsoft.com/en-us/azure/key-vault/keys/about-keys-details
// More info about the issue: http://coders-errand.com/malleability-ecdsa-signatures/
// More info about the fix: https://en.bitcoin.it/wiki/BIP_0062
func malleabilityECDSASignature(signature []byte) []byte {
	l := len(signature)
	hl := l / 2

	R := new(big.Int).SetBytes(signature[:hl])
	S := new(big.Int).SetBytes(signature[hl:l])
	if S.Cmp(secp256k1halfN) <= 0 {
		return signature
	}

	S2 := new(big.Int).Sub(secp256k1N, S)
	return append(R.Bytes(), S2.Bytes()...)
}
