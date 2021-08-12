package postgres

import (
	"context"
	"time"

	"github.com/consensys/quorum-key-manager/src/stores/database/models"
	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
	"github.com/consensys/quorum-key-manager/src/stores/database"

	"github.com/consensys/quorum-key-manager/pkg/errors"
)

type ETH1Accounts struct {
	storeID string
	logger  log.Logger
	client  postgres.Client
}

var _ database.ETH1Accounts = &ETH1Accounts{}

func NewETH1Accounts(storeID string, db postgres.Client, logger log.Logger) *ETH1Accounts {
	return &ETH1Accounts{
		storeID: storeID,
		logger:  logger,
		client:  db,
	}
}

func (ea ETH1Accounts) RunInTransaction(ctx context.Context, persist func(dbtx database.ETH1Accounts) error) error {
	return ea.client.RunInTransaction(ctx, func(dbTx postgres.Client) error {
		ea.client = dbTx
		return persist(&ea)
	})
}

func (ea *ETH1Accounts) Get(ctx context.Context, addr string) (*entities.ETH1Account, error) {
	eth1Acc := &models.ETH1Account{Address: addr, StoreID: ea.storeID}

	err := ea.client.SelectPK(ctx, eth1Acc)
	if err != nil {
		errMessage := "failed to get account"
		ea.logger.With("address", addr).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return eth1Acc.ToEntity(), nil
}

func (ea *ETH1Accounts) GetDeleted(ctx context.Context, addr string) (*entities.ETH1Account, error) {
	eth1Acc := &models.ETH1Account{Address: addr, StoreID: ea.storeID}

	err := ea.client.SelectDeletedPK(ctx, eth1Acc)
	if err != nil {
		errMessage := "failed to get deleted account"
		ea.logger.With("address", addr).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return eth1Acc.ToEntity(), nil
}

func (ea *ETH1Accounts) GetAll(ctx context.Context) ([]*entities.ETH1Account, error) {
	var eth1Accs []*models.ETH1Account

	err := ea.client.SelectWhere(ctx, &eth1Accs, "store_id = ?", ea.storeID)
	if err != nil {
		errMessage := "failed to get all accounts"
		ea.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	var accounts []*entities.ETH1Account
	for _, acc := range eth1Accs {
		accounts = append(accounts, acc.ToEntity())
	}

	return accounts, nil
}

func (ea *ETH1Accounts) GetAllDeleted(ctx context.Context) ([]*entities.ETH1Account, error) {
	var eth1Accs []*models.ETH1Account

	err := ea.client.SelectDeletedWhere(ctx, &eth1Accs, "store_id = ?", ea.storeID)
	if err != nil {
		errMessage := "failed to get all deleted accounts"
		ea.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	var accounts []*entities.ETH1Account
	for _, acc := range eth1Accs {
		accounts = append(accounts, acc.ToEntity())
	}

	return accounts, nil
}

func (ea *ETH1Accounts) Add(ctx context.Context, account *entities.ETH1Account) (*entities.ETH1Account, error) {
	accModel := models.NewETH1Account(account)
	accModel.StoreID = ea.storeID

	err := ea.client.Insert(ctx, accModel)
	if err != nil {
		errMessage := "failed to add account"
		ea.logger.With("address", account.Address).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return accModel.ToEntity(), nil
}

func (ea *ETH1Accounts) Update(ctx context.Context, account *entities.ETH1Account) (*entities.ETH1Account, error) {
	accModel := models.NewETH1Account(account)
	accModel.StoreID = ea.storeID

	err := ea.client.UpdatePK(ctx, accModel)
	if err != nil {
		errMessage := "failed to update account"
		ea.logger.With("address", account.Address).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return accModel.ToEntity(), nil
}

func (ea *ETH1Accounts) Delete(ctx context.Context, addr string) error {
	err := ea.client.DeletePK(ctx, &models.ETH1Account{Address: addr, StoreID: ea.storeID})
	if err != nil {
		errMessage := "failed to delete account"
		ea.logger.With("address", addr).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (ea *ETH1Accounts) Restore(ctx context.Context, account *entities.ETH1Account) error {
	accModel := models.NewETH1Account(account)
	accModel.StoreID = ea.storeID
	accModel.DeletedAt = time.Time{}
	err := ea.client.UndeletePK(ctx, accModel)
	if err != nil {
		errMessage := "failed to restore account"
		ea.logger.With("address", account.Address).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (ea *ETH1Accounts) Purge(ctx context.Context, addr string) error {
	err := ea.client.ForceDeletePK(ctx, &models.ETH1Account{Address: addr, StoreID: ea.storeID})
	if err != nil {
		errMessage := "failed to permanently delete account"
		ea.logger.With("address", addr).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}
