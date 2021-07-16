package connectors

import (
	"context"
	"math/big"

	"github.com/consensys/quorum-key-manager/pkg/ethereum"
	"github.com/consensys/quorum-key-manager/src/auth/policy"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/eth1"
	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/signer/core"
)

type Eth1Connector struct {
	store    eth1.Store
	logger   log.Logger
	resolver *policy.Resolver
}

var _ eth1.Store = Eth1Connector{}

func NewEth1Connector(store eth1.Store, resolvr *policy.Resolver, logger log.Logger) *Eth1Connector {
	return &Eth1Connector{
		store:    store,
		logger:   logger,
		resolver: resolvr,
	}
}

func (c Eth1Connector) Info(ctx context.Context) (*entities.StoreInfo, error) {
	result, err := c.store.Info(ctx)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("ethereum store info retrieved successfully")
	return result, nil
}

func (c Eth1Connector) Create(ctx context.Context, id string, attr *entities.Attributes) (*entities.ETH1Account, error) {
	logger := c.logger.With("id", id)
	logger.Debug("creating ethereum account")

	result, err := c.store.Create(ctx, id, attr)
	if err != nil {
		return nil, err
	}

	logger.Info("ethereum account created successfully")
	return result, nil
}

func (c Eth1Connector) Import(ctx context.Context, id string, privKey []byte, attr *entities.Attributes) (*entities.ETH1Account, error) {
	logger := c.logger.With("id", id)
	logger.Debug("importing ethereum account")

	result, err := c.store.Import(ctx, id, privKey, attr)
	if err != nil {
		return nil, err
	}

	logger.Info("ethereum account imported successfully")
	return result, nil
}

func (c Eth1Connector) Get(ctx context.Context, addr string) (*entities.ETH1Account, error) {
	logger := c.logger.With("address", addr)

	result, err := c.store.Get(ctx, addr)
	if err != nil {
		return nil, err
	}

	logger.Debug("ethereum account retrieved successfully")
	return result, nil
}

func (c Eth1Connector) GetAll(ctx context.Context) ([]*entities.ETH1Account, error) {
	result, err := c.store.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("all ethereum accounts retrieved successfully")
	return result, nil
}

func (c Eth1Connector) List(ctx context.Context) ([]string, error) {
	result, err := c.store.List(ctx)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("ethereum accounts listed successfully")
	return result, nil
}

func (c Eth1Connector) Update(ctx context.Context, addr string, attr *entities.Attributes) (*entities.ETH1Account, error) {
	logger := c.logger.With("address", addr)
	logger.Debug("updating ethereum account")

	result, err := c.store.Update(ctx, addr, attr)
	if err != nil {
		return nil, err
	}

	logger.Info("ethereum account updated successfully")
	return result, nil
}

func (c Eth1Connector) Delete(ctx context.Context, addr string) error {
	logger := c.logger.With("address", addr)
	logger.Debug("deleting ethereum account")

	err := c.store.Delete(ctx, addr)
	if err != nil {
		return err
	}

	logger.Info("ethereum account deleted successfully")
	return nil
}

func (c Eth1Connector) GetDeleted(ctx context.Context, addr string) (*entities.ETH1Account, error) {
	logger := c.logger.With("address", addr)

	result, err := c.store.GetDeleted(ctx, addr)
	if err != nil {
		return nil, err
	}

	logger.Debug("deleted ethereum account retrieved successfully")
	return result, nil
}

func (c Eth1Connector) ListDeleted(ctx context.Context) ([]string, error) {
	result, err := c.store.ListDeleted(ctx)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("deleted ethereum accounts listed successfully")
	return result, nil
}

func (c Eth1Connector) Undelete(ctx context.Context, addr string) error {
	logger := c.logger.With("address", addr)
	logger.Debug("restoring ethereum account")

	err := c.store.Undelete(ctx, addr)
	if err != nil {
		return err
	}

	logger.Info("ethereum account restored successfully")
	return nil
}

func (c Eth1Connector) Destroy(ctx context.Context, addr string) error {
	logger := c.logger.With("address", addr)
	logger.Debug("destroying ethereum account")

	err := c.store.Destroy(ctx, addr)
	if err != nil {
		return err
	}

	logger.Info("ethereum account was permanently deleted")
	return nil
}

func (c Eth1Connector) Sign(ctx context.Context, addr string, data []byte) ([]byte, error) {
	logger := c.logger.With("address", addr)

	result, err := c.store.Sign(ctx, addr, data)
	if err != nil {
		return nil, err
	}

	logger.Debug("payload signed successfully")
	return result, nil
}

func (c Eth1Connector) SignTypedData(ctx context.Context, addr string, typedData *core.TypedData) ([]byte, error) {
	logger := c.logger.With("address", addr)

	result, err := c.store.SignTypedData(ctx, addr, typedData)
	if err != nil {
		return nil, err
	}

	logger.Debug("typed data signed successfully")
	return result, nil
}

func (c Eth1Connector) SignTransaction(ctx context.Context, addr string, chainID *big.Int, tx *types.Transaction) ([]byte, error) {
	logger := c.logger.With("address", addr)
	result, err := c.store.SignTransaction(ctx, addr, chainID, tx)
	if err != nil {
		return nil, err
	}

	logger.Debug("transaction signed successfully")
	return result, nil
}

func (c Eth1Connector) SignEEA(ctx context.Context, addr string, chainID *big.Int, tx *types.Transaction, args *ethereum.PrivateArgs) ([]byte, error) {
	logger := c.logger.With("address", addr)
	result, err := c.store.SignEEA(ctx, addr, chainID, tx, args)
	if err != nil {
		return nil, err
	}

	logger.Debug("EEA transaction signed successfully")
	return result, nil
}

func (c Eth1Connector) SignPrivate(ctx context.Context, addr string, tx *quorumtypes.Transaction) ([]byte, error) {
	logger := c.logger.With("address", addr)
	result, err := c.store.SignPrivate(ctx, addr, tx)
	if err != nil {
		return nil, err
	}

	logger.Debug("private transaction signed successfully")
	return result, nil
}

func (c Eth1Connector) ECRevocer(ctx context.Context, data, sig []byte) (string, error) {
	result, err := c.store.ECRevocer(ctx, data, sig)
	if err != nil {
		return "", err
	}

	c.logger.Debug("EC recovered successfully")
	return result, nil
}

func (c Eth1Connector) Verify(ctx context.Context, addr string, data, sig []byte) error {
	err := c.store.Verify(ctx, addr, data, sig)
	if err != nil {
		return err
	}

	c.logger.Debug("data verified successfully")
	return nil
}

func (c Eth1Connector) VerifyTypedData(ctx context.Context, addr string, typedData *core.TypedData, sig []byte) error {
	err := c.store.VerifyTypedData(ctx, addr, typedData, sig)
	if err != nil {
		c.logger.WithError(err).Error("failed to verify typed data")
		return err
	}

	c.logger.Debug("typed data verified successfully")
	return nil
}

func (c Eth1Connector) Encrypt(ctx context.Context, addr string, data []byte) ([]byte, error) {
	logger := c.logger.With("address", addr)
	result, err := c.store.Encrypt(ctx, addr, data)
	if err != nil {
		c.logger.WithError(err).Error("failed to encrypt data")
		return nil, err
	}

	logger.Debug("data encrypted successfully")
	return result, nil
}

func (c Eth1Connector) Decrypt(ctx context.Context, addr string, data []byte) ([]byte, error) {
	logger := c.logger.With("address", addr)
	result, err := c.store.Decrypt(ctx, addr, data)
	if err != nil {
		c.logger.WithError(err).Error("failed to decrypt data")
		return nil, err
	}

	logger.Debug("data decrypted successfully")
	return result, nil
}
