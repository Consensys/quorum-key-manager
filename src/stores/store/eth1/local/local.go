package local

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/consensys/quorum-key-manager/src/stores/connectors"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/ethereum"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"github.com/consensys/quorum-key-manager/src/stores/store/database"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/eth1"
	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/signer/core"
	"golang.org/x/crypto/sha3"
)

// Values copied from github.com/ethereum/go-ethereum/crypto/crypto.go
var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

var eth1KeyAlgo = &entities.Algorithm{
	Type:          entities.Ecdsa,
	EllipticCurve: entities.Secp256k1,
}

type Store struct {
	keysConnector connectors.KeysConnector
	db            database.Database
	logger        log.Logger
}

var _ eth1.Store = &Store{}

func New(keysConnector connectors.KeysConnector, db database.Database, logger log.Logger) *Store {
	return &Store{
		keysConnector: keysConnector,
		logger:        logger,
		db:            db,
	}
}

func (s *Store) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Create(ctx context.Context, id string, attr *entities.Attributes) (*entities.ETH1Account, error) {
	key, err := s.keysConnector.Create(ctx, id, eth1KeyAlgo, attr)
	if err != nil {
		return nil, err
	}

	return s.db.ETH1Accounts().Add(ctx, parseKey(key, attr))
}

func (s *Store) Import(ctx context.Context, id string, privKey []byte, attr *entities.Attributes) (*entities.ETH1Account, error) {
	key, err := s.keysConnector.Import(ctx, id, privKey, eth1KeyAlgo, attr)
	if err != nil {
		return nil, err
	}

	return s.db.ETH1Accounts().Add(ctx, parseKey(key, attr))
}

func (s *Store) Get(ctx context.Context, addr string) (*entities.ETH1Account, error) {
	return s.db.ETH1Accounts().Get(ctx, addr)
}

func (s *Store) GetAll(ctx context.Context) ([]*entities.ETH1Account, error) {
	return s.db.ETH1Accounts().GetAll(ctx)
}

func (s *Store) List(ctx context.Context) ([]string, error) {
	addresses := []string{}
	accountsRetrieved, err := s.db.ETH1Accounts().GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, acc := range accountsRetrieved {
		addresses = append(addresses, acc.Address.Hex())
	}

	return addresses, nil
}

func (s *Store) Update(ctx context.Context, addr string, attr *entities.Attributes) (*entities.ETH1Account, error) {
	account, err := s.db.ETH1Accounts().Get(ctx, addr)
	if err != nil {
		return nil, err
	}
	account.Tags = attr.Tags

	return s.db.ETH1Accounts().Update(ctx, account)
}

func (s *Store) Delete(ctx context.Context, addr string) error {
	acc, err := s.db.ETH1Accounts().Get(ctx, addr)
	if err != nil {
		return err
	}

	return s.db.RunInTransaction(ctx, func(dbtx database.Database) error {
		err = s.db.ETH1Accounts().Delete(ctx, addr)
		if err != nil {
			return err
		}

		err = s.keysConnector.Delete(ctx, acc.KeyID)
		if err != nil && !errors.IsNotSupportedError(err) { // If the underlying store does not support deleting, we only delete in DB
			return err
		}

		return nil
	})
}

func (s *Store) GetDeleted(ctx context.Context, addr string) (*entities.ETH1Account, error) {
	return s.db.ETH1Accounts().GetDeleted(ctx, addr)
}

func (s *Store) ListDeleted(ctx context.Context) ([]string, error) {
	addresses := []string{}
	accountsRetrieved, err := s.db.ETH1Accounts().GetAllDeleted(ctx)
	if err != nil {
		return nil, err
	}

	for _, acc := range accountsRetrieved {
		addresses = append(addresses, acc.Address.Hex())
	}

	return addresses, nil
}

func (s *Store) Undelete(ctx context.Context, addr string) error {
	account, err := s.db.ETH1Accounts().GetDeleted(ctx, addr)
	if err != nil {
		return err
	}

	return s.db.RunInTransaction(ctx, func(dbtx database.Database) error {
		derr := s.db.ETH1Accounts().Restore(ctx, account)
		if derr != nil {
			return derr
		}

		return s.keysConnector.Restore(ctx, account.KeyID)
	})
}

func (s *Store) Destroy(ctx context.Context, addr string) error {
	account, err := s.db.ETH1Accounts().GetDeleted(ctx, addr)
	if err != nil {
		return err
	}

	return s.db.RunInTransaction(ctx, func(dbtx database.Database) error {
		derr := s.db.ETH1Accounts().Purge(ctx, addr)
		if derr != nil {
			return derr
		}

		return s.keysConnector.Destroy(ctx, account.KeyID)
	})
}

func (s *Store) Sign(ctx context.Context, addr string, data []byte) ([]byte, error) {
	return s.SignData(ctx, addr, crypto.Keccak256(data))
}

func (s *Store) SignTypedData(ctx context.Context, addr string, typedData *core.TypedData) ([]byte, error) {
	encodedData, err := getEIP712EncodedData(typedData)
	if err != nil {
		return nil, err
	}

	return s.Sign(ctx, addr, []byte(encodedData))
}

func (s *Store) SignTransaction(ctx context.Context, addr string, chainID *big.Int, tx *types.Transaction) ([]byte, error) {
	signer := types.NewEIP155Signer(chainID)
	txData := signer.Hash(tx).Bytes()
	signature, err := s.SignData(ctx, addr, txData)
	if err != nil {
		return nil, err
	}

	signedTx, err := tx.WithSignature(signer, signature)
	if err != nil {
		errMessage := "failed to set transaction signature"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.DependencyFailureError(errMessage)
	}

	signedRaw, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		errMessage := "failed to RLP encode signed transaction"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.EncodingError(errMessage)
	}

	return signedRaw, nil
}

func (s *Store) SignEEA(ctx context.Context, addr string, chainID *big.Int, tx *types.Transaction, args *ethereum.PrivateArgs) ([]byte, error) {
	privateFromEncoded, err := base64.StdEncoding.DecodeString(*args.PrivateFrom)
	if err != nil {
		errMessage := "invalid privateFrom param"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	privateRecipientEncoded, err := getEncodedPrivateRecipient(args.PrivacyGroupID, args.PrivateFor)
	if err != nil {
		errMessage := "invalid privacyGroupID or privateFor"
		s.logger.WithError(err).Error(errMessage)
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
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	signature, err := s.SignData(ctx, addr, hash[:])
	if err != nil {
		return nil, err
	}

	signedTx, err := tx.WithSignature(types.NewEIP155Signer(chainID), signature)
	if err != nil {
		errMessage := "failed to set eea transaction signature"
		s.logger.WithError(err).Error(errMessage)
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
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.EncodingError(errMessage)
	}

	return signedRaw, nil
}

func (s *Store) SignPrivate(ctx context.Context, addr string, tx *quorumtypes.Transaction) ([]byte, error) {
	signer := quorumtypes.QuorumPrivateTxSigner{}
	txData := signer.Hash(tx).Bytes()
	signature, err := s.SignData(ctx, addr, txData)
	if err != nil {
		return nil, err
	}

	signedTx, err := tx.WithSignature(signer, signature)
	if err != nil {
		errMessage := "failed to set quorum private transaction signature"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.DependencyFailureError(errMessage)
	}

	signedRaw, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		errMessage := "failed to RLP encode signed quorum private transaction"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.EncodingError(errMessage)
	}

	return signedRaw, nil
}

func (s *Store) ECRevocer(_ context.Context, data, sig []byte) (string, error) {
	pubKey, err := crypto.SigToPub(crypto.Keccak256(data), sig)
	if err != nil {
		errMessage := "failed to recover public key, please verify your signature and payload"
		s.logger.WithError(err).Error(errMessage)
		return "", errors.InvalidParameterError(errMessage)
	}

	return crypto.PubkeyToAddress(*pubKey).Hex(), nil
}

func (s *Store) Verify(ctx context.Context, addr string, data, sig []byte) error {
	recoveredAddress, err := s.ECRevocer(ctx, data, sig)
	if err != nil {
		return err
	}

	if addr != recoveredAddress {
		errMessage := "failed to verify signature: recovered address does not match the expected one or payload is malformed"
		s.logger.WithError(err).Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	return nil
}

func (s *Store) VerifyTypedData(ctx context.Context, addr string, typedData *core.TypedData, sig []byte) error {
	encodedData, err := getEIP712EncodedData(typedData)
	if err != nil {
		errMessage := "failed to generate EIP-712 encoded data"
		s.logger.WithError(err).Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	return s.Verify(ctx, addr, []byte(encodedData), sig)
}

func (s *Store) Encrypt(ctx context.Context, addr string, data []byte) ([]byte, error) {
	account, err := s.db.ETH1Accounts().Get(ctx, addr)
	if err != nil {
		return nil, err
	}

	return s.keysConnector.Encrypt(ctx, account.KeyID, data)
}

func (s *Store) Decrypt(ctx context.Context, addr string, data []byte) ([]byte, error) {
	account, err := s.db.ETH1Accounts().Get(ctx, addr)
	if err != nil {
		return nil, err
	}

	return s.keysConnector.Decrypt(ctx, account.KeyID, data)
}

func getEIP712EncodedData(typedData *core.TypedData) (string, error) {
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return "", err
	}

	domainSeparatorHash, err := typedData.HashStruct(formatters.EIP712DomainLabel, typedData.Domain.Map())
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("\x19\x01%s%s", domainSeparatorHash, typedDataHash), nil
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

func eeaHash(object interface{}) (hash common.Hash, err error) {
	hashAlgo := sha3.NewLegacyKeccak256()
	err = rlp.Encode(hashAlgo, object)
	if err != nil {
		return common.Hash{}, err
	}

	hashAlgo.Sum(hash[:0])

	return hash, nil
}

func (s *Store) SignData(ctx context.Context, addr string, data []byte) ([]byte, error) {
	account, err := s.Get(ctx, addr)
	if err != nil {
		return nil, err
	}

	signature, err := s.keysConnector.Sign(ctx, account.KeyID, data)
	if err != nil {
		return nil, err
	}

	// Recover the recID, please read: http://coders-errand.com/ecrecover-signature-verification-ethereum/
	for _, recID := range []byte{0, 1} {
		appendedSignature := append(malleabilityECDSASignature(signature), recID)
		recoveredPubKey, err := crypto.SigToPub(data, appendedSignature)
		if err != nil {
			errMessage := "failed to recover public key candidate with appended recID"
			s.logger.WithError(err).Error(errMessage, "recID", recID)
			return nil, errors.InvalidParameterError(errMessage)
		}

		if bytes.Equal(crypto.FromECDSAPub(recoveredPubKey), account.PublicKey) {
			return appendedSignature, nil
		}
	}

	errMessage := "failed to compute recovery ID"
	s.logger.Error(errMessage)
	return nil, errors.DependencyFailureError(errMessage)
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
