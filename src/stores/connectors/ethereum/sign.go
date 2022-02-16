package eth

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"

	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"

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
	secp256k1halfN, _ = new(big.Int).SetString("7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0", 16)
	maxRetries        = 10
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

	signature, err := c.signHomestead(ctx, addr, crypto.Keccak256(ethereum.GetEIP191EncodedData(data)))
	if err != nil {
		return nil, err
	}

	logger.Debug("message signed successfully (eip-191)")
	return signature, nil
}

func (c Connector) SignTypedData(ctx context.Context, addr common.Address, typedData *core.TypedData) ([]byte, error) {
	logger := c.logger.With("address", addr.Hex())

	encodedData, err := ethereum.GetEIP712EncodedData(typedData)
	if err != nil {
		errMessage := "failed to format typed data"
		logger.WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	signature, err := c.signHomestead(ctx, addr, crypto.Keccak256(encodedData))
	if err != nil {
		return nil, err
	}

	logger.Debug("typed data signed successfully (eip-712)")
	return signature, nil
}

func (c Connector) SignTransaction(ctx context.Context, addr common.Address, chainID *big.Int, tx *types.Transaction) ([]byte, error) {
	logger := c.logger.With("address", addr.Hex())

	signer := types.NewLondonSigner(chainID)
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

	signedRaw, err := signedTx.MarshalBinary()
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
		errMessage := "invalid 'privateFrom'"
		c.logger.WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	privateRecipientEncoded, err := getEncodedPrivateRecipient(args.PrivacyGroupID, args.PrivateFor)
	if err != nil {
		errMessage := "invalid 'privacyGroupID' or 'privateFor'"
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

	signature, err := c.sign(ctx, addr, hash.Bytes())
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
	err := c.authorizator.CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount})
	if err != nil {
		return nil, err
	}

	acc, err := c.db.Get(ctx, addr.Hex())
	if err != nil {
		return nil, err
	}

	var signature []byte
	var retry int
	for retry = maxRetries; retry > 0; retry-- {
		signature, err = c.store.Sign(ctx, acc.KeyID, data, ethAlgo)
		if err != nil {
			return nil, err
		}

		// If we get a malleable signature, we retry
		if !isMalleableECDSASignature(signature) {
			break
		}

		c.logger.Debug("malleable signature retrieved, retryng", "signature", hexutil.Encode(signature))
	}

	if retry == 0 {
		errMessage := "failed to generate a non malleable signature"
		c.logger.Error(errMessage)
		return nil, errors.DependencyFailureError(errMessage)
	}

	// Recover the recID, please read: http://coders-errand.com/ecrecover-signature-verification-ethereum/
	for _, recID := range []byte{0, 1} {
		appendedSignature := append(signature, recID)
		var recoveredPubKey *ecdsa.PublicKey
		recoveredPubKey, err = crypto.SigToPub(data, appendedSignature)
		if err != nil {
			errMessage := "failed to recover public key candidate with appended recID"
			c.logger.WithError(err).Error(errMessage, "recID", recID)
			return nil, errors.CryptoOperationError(errMessage)
		}

		if bytes.Equal(crypto.FromECDSAPub(recoveredPubKey), acc.PublicKey) {
			return appendedSignature, nil
		}
	}

	errMessage := "failed to recover public key candidate"
	c.logger.Error(errMessage)
	return nil, errors.CryptoOperationError(errMessage)
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

func (c Connector) signHomestead(ctx context.Context, addr common.Address, data []byte) ([]byte, error) {
	signature, err := c.sign(ctx, addr, data)
	if err != nil {
		return nil, err
	}

	signature[crypto.RecoveryIDOffset] += 27

	return signature, nil
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
func isMalleableECDSASignature(signature []byte) bool {
	return new(big.Int).SetBytes(signature[32:]).Cmp(secp256k1halfN) >= 0
}
