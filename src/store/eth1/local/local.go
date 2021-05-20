package local

import (
	"context"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/eth1"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core"
	"math/big"
	"sync"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

const domainLabel = "EIP712Domain"

var eth1KeyAlgo = &entities.Algorithm{
	Type:          entities.Ecdsa,
	EllipticCurve: entities.Secp256k1,
}

// Store is an implementation of ethereum (ETH1) store relying on an underlying key store
type Store struct {
	keyStore        keys.Store
	addrToId        map[string]string
	deletedAddrToId map[string]string
	mux             sync.RWMutex
}

var _ eth1.Store = &Store{}

// New creates an HashiCorp key store
func New(keyStore keys.Store) *Store {
	return &Store{
		mux:             sync.RWMutex{},
		keyStore:        keyStore,
		addrToId:        make(map[string]string),
		deletedAddrToId: make(map[string]string),
	}
}

func (s *Store) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

// Create an Ethereum account
func (s *Store) Create(ctx context.Context, id string, attr *entities.Attributes) (*entities.ETH1Account, error) {
	key, err := s.keyStore.Create(ctx, id, eth1KeyAlgo, attr)
	if err != nil {
		return nil, err
	}

	acc, err := parseKey(key)
	if err != nil {
		return nil, err
	}

	s.addID(acc.Address, acc.ID)
	return acc, nil
}

// Import an ETH1 account
func (s *Store) Import(ctx context.Context, id string, privKey []byte, attr *entities.Attributes) (*entities.ETH1Account, error) {
	key, err := s.keyStore.Import(ctx, id, privKey, eth1KeyAlgo, attr)
	if err != nil {
		return nil, err
	}

	acc, err := parseKey(key)
	if err != nil {
		return nil, err
	}

	s.addID(acc.Address, acc.ID)
	return acc, nil
}

// Get an account
func (s *Store) Get(ctx context.Context, addr string) (*entities.ETH1Account, error) {
	id, err := s.getID(addr)
	if err != nil {
		return nil, err
	}

	key, err := s.keyStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return parseKey(key)
}

// Get all accounts
func (s *Store) GetAll(ctx context.Context) ([]*entities.ETH1Account, error) {
	var accounts = make([]*entities.ETH1Account, len(s.addrToId))

	for _, id := range s.addrToId {
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

// Get all account ids
func (s *Store) List(_ context.Context) ([]string, error) {
	addresses := make([]string, len(s.addrToId))

	for address, _ := range s.addrToId {
		addresses = append(addresses, address)
	}

	return addresses, nil
}

// Update account tags
func (s *Store) Update(ctx context.Context, addr string, attr *entities.Attributes) (*entities.ETH1Account, error) {
	id, err := s.getID(addr)
	if err != nil {
		return nil, err
	}

	key, err := s.keyStore.Update(ctx, id, attr)
	if err != nil {
		return nil, err
	}

	return parseKey(key)
}

// Delete an account
func (s *Store) Delete(ctx context.Context, addr string) error {
	id, err := s.getID(addr)
	if err != nil {
		return err
	}

	err = s.keyStore.Delete(ctx, id)
	if err != nil {
		return err
	}

	s.removeID(addr)
	s.addDeletedID(addr, id)

	return nil
}

// Gets a deleted account
func (s *Store) GetDeleted(ctx context.Context, addr string) (*entities.ETH1Account, error) {
	id, err := s.getDeletedID(addr)
	if err != nil {
		return nil, err
	}

	key, err := s.keyStore.GetDeleted(ctx, id)
	if err != nil {
		return nil, err
	}

	return parseKey(key)
}

// Lists all deleted accounts
func (s *Store) ListDeleted(_ context.Context) ([]string, error) {
	addresses := make([]string, len(s.deletedAddrToId))

	for addr, _ := range s.deletedAddrToId {
		addresses = append(addresses, addr)
	}

	return addresses, nil
}

// Undelete a previously deleted account
func (s *Store) Undelete(ctx context.Context, addr string) error {
	id, err := s.getDeletedID(addr)
	if err != nil {
		return err
	}

	err = s.keyStore.Undelete(ctx, id)
	if err != nil {
		return err
	}

	s.removeDeletedID(addr)
	s.addID(addr, id)

	return nil
}

// Destroy an account permanently
func (s *Store) Destroy(ctx context.Context, addr string) error {
	id, err := s.getDeletedID(addr)
	if err != nil {
		return err
	}

	err = s.keyStore.Destroy(ctx, id)
	if err != nil {
		return err
	}

	s.removeDeletedID(addr)
	return nil
}

// Sign any arbitrary data
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

// Sign EIP-712 formatted data using the specified account
func (s *Store) SignTypedData(ctx context.Context, addr string, typedData *core.TypedData) ([]byte, error) {
	encodedData, err := getEIP712EncodedData(typedData)
	if err != nil {
		return nil, err
	}

	key, err := s.Get(ctx, addr)
	if err != nil {
		return nil, err
	}

	signature, err := s.Sign(ctx, addr, []byte(encodedData))
	if err != nil {
		return nil, err
	}

	return appendRecID(signature, key.PublicKey)
}

func (s *Store) SignTransaction(ctx context.Context, addr string, chainID *big.Int, tx *ethereum.TxData) ([]byte, error) {
	key, err := s.Get(ctx, addr)
	if err != nil {
		return nil, err
	}

	ethTx := types.NewTransaction(tx.Nonce, *tx.To, tx.Value, tx.GasLimit, tx.GasPrice, tx.Data)
	signature, err := s.Sign(ctx, addr, types.NewEIP155Signer(chainID).Hash(ethTx).Bytes())
	if err != nil {
		return nil, err
	}

	// The signature is [R||S] and we want to add the V value to make it compatible with Ethereum [R||S||V]
	signature, err = appendRecID(signature, key.PublicKey)
	if err != nil {
		return nil, err
	}

	ethSignature, err := ethereum.EIP155Signature(signature, chainID)
	if err != nil {
		return nil, errors.EncodingError("failed to recover signature V value")
	}

	return ethSignature, nil
}

// SignEEA transaction
func (s *Store) SignEEA(ctx context.Context, addr string, chainID *big.Int, tx *ethereum.EEATxData, args *ethereum.PrivateArgs) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

// SignPrivate transaction
func (s *Store) SignPrivate(ctx context.Context, addr string, tx *ethereum.TxData) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

// ECRevocer returns the address from a signature and data
func (s *Store) ECRevocer(_ context.Context, sig, data []byte) (string, error) {
	pubKey, err := crypto.SigToPub(crypto.Keccak256(data), sig)
	if err != nil {
		return "", errors.InvalidParameterError("failed to recover public key, please verify your signature and payload")
	}

	return crypto.PubkeyToAddress(*pubKey).Hex(), nil
}

// Verify verifies that a signature belongs to a given address
func (s *Store) Verify(ctx context.Context, addr string, sig, payload []byte) error {
	recoveredAddress, err := s.ECRevocer(ctx, payload, sig)
	if err != nil {
		return err
	}

	if addr != recoveredAddress {
		return errors.InvalidParameterError("failed to verify signature: recovered address does not match the expected one or payload is malformed")
	}

	return nil
}

// Verify verifies that a typed data signature belongs to a given address
func (s *Store) VerifyTypedData(ctx context.Context, addr string, sig []byte, typedData *core.TypedData) error {
	encodedData, err := getEIP712EncodedData(typedData)
	if err != nil {
		return err
	}

	return s.Verify(ctx, addr, sig, []byte(encodedData))
}

// Encrypt any arbitrary data using a specified account
func (s *Store) Encrypt(ctx context.Context, addr string, data []byte) ([]byte, error) {
	id, err := s.getID(addr)
	if err != nil {
		return nil, err
	}

	return s.keyStore.Encrypt(ctx, id, data)
}

// Decrypt a single block of encrypted data.
func (s *Store) Decrypt(ctx context.Context, addr string, data []byte) ([]byte, error) {
	id, err := s.getID(addr)
	if err != nil {
		return nil, err
	}

	return s.keyStore.Decrypt(ctx, id, data)
}

func (s *Store) getID(addr string) (string, error) {
	id, ok := s.addrToId[addr]
	if !ok {
		return "", errors.NotFoundError("account %s was not found", addr)
	}

	return id, nil
}

func (s *Store) getDeletedID(addr string) (string, error) {
	id, ok := s.deletedAddrToId[addr]
	if !ok {
		return "", errors.NotFoundError("deleted account %s was not found", addr)
	}

	return id, nil
}

func (s *Store) addID(addr, id string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.addrToId[addr] = id
}

func (s *Store) addDeletedID(addr, id string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.deletedAddrToId[addr] = id
}

func (s *Store) removeID(addr string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.addrToId, addr)
}

func (s *Store) removeDeletedID(addr string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.deletedAddrToId, addr)
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
