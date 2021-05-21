package local

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/database"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/eth1"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

const domainLabel = "EIP712Domain"

var eth1KeyAlgo = &entities.Algorithm{
	Type:          entities.Ecdsa,
	EllipticCurve: entities.Secp256k1,
}

type Store struct {
	keyStore       keys.Store
	eth1AccountsDB database.Database
	logger         *log.Logger
}

var _ eth1.Store = &Store{}

func New(keyStore keys.Store, eth1AccountsDB database.Database, logger *log.Logger) *Store {
	return &Store{
		keyStore:       keyStore,
		logger:         logger,
		eth1AccountsDB: eth1AccountsDB,
	}
}

func (s *Store) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Create(ctx context.Context, id string, attr *entities.Attributes) (*entities.ETH1Account, error) {
	key, err := s.keyStore.Create(ctx, id, eth1KeyAlgo, attr)
	if err != nil {
		return nil, err
	}

	acc, err := parseKey(key)
	if err != nil {
		return nil, err
	}

	err = s.eth1AccountsDB.AddID(ctx, acc.Address, acc.ID)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *Store) Import(ctx context.Context, id string, privKey []byte, attr *entities.Attributes) (*entities.ETH1Account, error) {
	key, err := s.keyStore.Import(ctx, id, privKey, eth1KeyAlgo, attr)
	if err != nil {
		return nil, err
	}

	acc, err := parseKey(key)
	if err != nil {
		return nil, err
	}

	err = s.eth1AccountsDB.AddID(ctx, acc.Address, acc.ID)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *Store) Get(ctx context.Context, addr string) (*entities.ETH1Account, error) {
	id, err := s.eth1AccountsDB.GetID(ctx, addr)
	if err != nil {
		return nil, err
	}

	key, err := s.keyStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return parseKey(key)
}

func (s *Store) GetAll(ctx context.Context) ([]*entities.ETH1Account, error) {
	ids, err := s.eth1AccountsDB.GetAllIDs(ctx)
	if err != nil {
		return nil, err
	}

	var accounts = make([]*entities.ETH1Account, len(ids))
	for _, id := range ids {
		key, err := s.keyStore.Get(ctx, id)
		if err != nil {
			return nil, err
		}

		account, err := parseKey(key)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *Store) List(ctx context.Context) ([]string, error) {
	return s.eth1AccountsDB.GetAll(ctx)
}

func (s *Store) Update(ctx context.Context, addr string, attr *entities.Attributes) (*entities.ETH1Account, error) {
	id, err := s.eth1AccountsDB.GetID(ctx, addr)
	if err != nil {
		return nil, err
	}

	key, err := s.keyStore.Update(ctx, id, attr)
	if err != nil {
		return nil, err
	}

	return parseKey(key)
}

func (s *Store) Delete(ctx context.Context, addr string) error {
	id, err := s.eth1AccountsDB.GetID(ctx, addr)
	if err != nil {
		return err
	}

	err = s.keyStore.Delete(ctx, id)
	if err != nil {
		return err
	}

	err = s.eth1AccountsDB.RemoveID(ctx, addr)
	if err != nil {
		return err
	}

	return s.eth1AccountsDB.AddDeletedID(ctx, addr, id)
}

func (s *Store) GetDeleted(ctx context.Context, addr string) (*entities.ETH1Account, error) {
	id, err := s.eth1AccountsDB.GetDeletedID(ctx, addr)
	if err != nil {
		return nil, err
	}

	key, err := s.keyStore.GetDeleted(ctx, id)
	if err != nil {
		return nil, err
	}

	return parseKey(key)
}

func (s *Store) ListDeleted(ctx context.Context) ([]string, error) {
	return s.eth1AccountsDB.GetAllDeleted(ctx)
}

func (s *Store) Undelete(ctx context.Context, addr string) error {
	id, err := s.eth1AccountsDB.GetDeletedID(ctx, addr)
	if err != nil {
		return err
	}

	err = s.keyStore.Undelete(ctx, id)
	if err != nil {
		return err
	}

	err = s.eth1AccountsDB.RemoveDeletedID(ctx, addr)
	if err != nil {
		return err
	}

	return s.eth1AccountsDB.AddID(ctx, addr, id)
}

func (s *Store) Destroy(ctx context.Context, addr string) error {
	id, err := s.eth1AccountsDB.GetDeletedID(ctx, addr)
	if err != nil {
		return err
	}

	err = s.keyStore.Destroy(ctx, id)
	if err != nil {
		return err
	}

	return s.eth1AccountsDB.RemoveDeletedID(ctx, addr)
}

func (s *Store) Sign(ctx context.Context, addr string, data []byte) ([]byte, error) {
	key, err := s.Get(ctx, addr)
	if err != nil {
		return nil, err
	}

	signature, err := s.keyStore.Sign(ctx, key.ID, data)
	if err != nil {
		return nil, err
	}

	return appendRecID(signature, key.PublicKey)
}

func (s *Store) SignTypedData(ctx context.Context, addr string, typedData *core.TypedData) ([]byte, error) {
	encodedData, err := getEIP712EncodedData(typedData)
	if err != nil {
		return nil, err
	}

	key, err := s.Get(ctx, addr)
	if err != nil {
		return nil, err
	}

	signature, err := s.Sign(ctx, addr, crypto.Keccak256([]byte(encodedData)))
	if err != nil {
		return nil, err
	}

	return appendRecID(signature, key.PublicKey)
}

func (s *Store) SignTransaction(ctx context.Context, addr string, chainID *big.Int, tx *types.Transaction) ([]byte, error) {
	logger := s.logger.WithField("address", addr)

	signer := types.NewEIP155Signer(chainID)
	signature, err := s.Sign(ctx, addr, signer.Hash(tx).Bytes())
	if err != nil {
		return nil, err
	}

	ethSignature, err := parseSignatureValues(tx, signature, signer)
	if err != nil {
		errMessage := "failed to generate transaction signature"
		logger.WithError(err).Error(errMessage)
		return nil, errors.EncodingError(errMessage)
	}

	return ethSignature, nil
}

func (s *Store) SignEEA(ctx context.Context, addr string, chainID *big.Int, tx *ethereum.EEATxData, args *ethereum.PrivateArgs) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) SignPrivate(ctx context.Context, addr string, tx *types.Transaction) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) ECRevocer(_ context.Context, data, sig []byte) (string, error) {
	pubKey, err := crypto.SigToPub(data, sig)
	if err != nil {
		s.logger.WithError(err).Error("failed to recover public key")
		return "", errors.InvalidParameterError("failed to recover public key, please verify your signature and payload")
	}

	return crypto.PubkeyToAddress(*pubKey).Hex(), nil
}

func (s *Store) Verify(ctx context.Context, addr string, data, sig []byte) error {
	recoveredAddress, err := s.ECRevocer(ctx, data, sig)
	if err != nil {
		return err
	}

	if addr != recoveredAddress {
		return errors.InvalidParameterError("failed to verify signature: recovered address does not match the expected one or payload is malformed")
	}

	return nil
}

func (s *Store) VerifyTypedData(ctx context.Context, addr string, sig []byte, typedData *core.TypedData) error {
	encodedData, err := getEIP712EncodedData(typedData)
	if err != nil {
		return err
	}

	return s.Verify(ctx, addr, sig, []byte(encodedData))
}

func (s *Store) Encrypt(ctx context.Context, addr string, data []byte) ([]byte, error) {
	id, err := s.eth1AccountsDB.GetID(ctx, addr)
	if err != nil {
		return nil, err
	}

	return s.keyStore.Encrypt(ctx, id, data)
}

func (s *Store) Decrypt(ctx context.Context, addr string, data []byte) ([]byte, error) {
	id, err := s.eth1AccountsDB.GetID(ctx, addr)
	if err != nil {
		return nil, err
	}

	return s.keyStore.Decrypt(ctx, id, data)
}

func getEIP712EncodedData(typedData *core.TypedData) (string, error) {
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return "", errors.InvalidParameterError("invalid typed data message")
	}

	domainSeparatorHash, err := typedData.HashStruct(domainLabel, typedData.Domain.Map())
	if err != nil {
		return "", errors.InvalidParameterError("invalid domain separator")
	}

	return fmt.Sprintf("\x19\x01%s%s", domainSeparatorHash, typedDataHash), nil
}

func appendRecID(sig, pubKey []byte) ([]byte, error) {
	recID, err := parseRecID(pubKey)
	if err != nil {
		return nil, err
	}

	return append(sig, *recID), nil
}

