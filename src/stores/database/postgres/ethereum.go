package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/consensys/quorum-key-manager/src/stores/database/models"
	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
	"github.com/consensys/quorum-key-manager/src/stores/database"

	"github.com/consensys/quorum-key-manager/pkg/errors"
)

type ETHAccounts struct {
	storeID string
	logger  log.Logger
	client  postgres.Client
}

var _ database.ETHAccounts = &ETHAccounts{}

func NewETHAccounts(storeID string, db postgres.Client, logger log.Logger) *ETHAccounts {
	return &ETHAccounts{
		storeID: storeID,
		logger:  logger,
		client:  db,
	}
}

func (ea ETHAccounts) RunInTransaction(ctx context.Context, persist func(dbtx database.ETHAccounts) error) error {
	return ea.client.RunInTransaction(ctx, func(dbTx postgres.Client) error {
		ea.client = dbTx
		return persist(&ea)
	})
}

func (ea *ETHAccounts) Get(ctx context.Context, addr string) (*entities.ETHAccount, error) {
	ethAcc := &models.ETHAccount{Address: addr, StoreID: ea.storeID}

	err := ea.client.SelectPK(ctx, ethAcc)
	if err != nil {
		errMessage := "failed to get account"
		ea.logger.With("address", addr).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return ethAcc.ToEntity(), nil
}

func (ea *ETHAccounts) GetDeleted(ctx context.Context, addr string) (*entities.ETHAccount, error) {
	ethAcc := &models.ETHAccount{Address: addr, StoreID: ea.storeID}

	err := ea.client.SelectDeletedPK(ctx, ethAcc)
	if err != nil {
		errMessage := "failed to get deleted account"
		ea.logger.With("address", addr).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return ethAcc.ToEntity(), nil
}

func (ea *ETHAccounts) GetAll(ctx context.Context) ([]*entities.ETHAccount, error) {
	var ethAccs []*models.ETHAccount

	err := ea.client.SelectWhere(ctx, &ethAccs, "store_id = ?", ea.storeID)
	if err != nil {
		errMessage := "failed to get all accounts"
		ea.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	var accounts []*entities.ETHAccount
	for _, acc := range ethAccs {
		accounts = append(accounts, acc.ToEntity())
	}

	return accounts, nil
}

func (ea *ETHAccounts) GetAllDeleted(ctx context.Context) ([]*entities.ETHAccount, error) {
	var ethAccs []*models.ETHAccount

	err := ea.client.SelectDeletedWhere(ctx, &ethAccs, "store_id = ?", ea.storeID)
	if err != nil {
		errMessage := "failed to get all deleted accounts"
		ea.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	var accounts []*entities.ETHAccount
	for _, acc := range ethAccs {
		accounts = append(accounts, acc.ToEntity())
	}

	return accounts, nil
}

func (ea *ETHAccounts) ListAddresses(ctx context.Context, isDeleted bool, limit, offset int) ([]string, error) {
	var ids = []string{}
	var err error
	var query string
	args := []interface{}{ea.storeID}

	switch {
	case limit != 0 || offset != 0:
		query = fmt.Sprintf("SELECT (array_agg(address ORDER BY created_at ASC))[%d:%d] FROM eth_accounts WHERE store_id = ?", offset+1, offset+limit)
	default:
		query = "SELECT array_agg(address ORDER BY created_at ASC) FROM eth_accounts WHERE store_id = ?"
	}

	if isDeleted {
		query = fmt.Sprintf("%s AND deleted_at is NOT NULL", query)
	} else {
		query = fmt.Sprintf("%s AND deleted_at is NULL", query)
	}

	err = ea.client.Query(ctx, &ids, query, args...)
	if err != nil {
		errMessage := "failed to list of ethereum addresses"
		ea.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return ids, nil
}

func (ea *ETHAccounts) Add(ctx context.Context, account *entities.ETHAccount) (*entities.ETHAccount, error) {
	accModel := models.NewETHAccount(account)
	accModel.StoreID = ea.storeID
	accModel.CreatedAt = time.Now()
	accModel.UpdatedAt = time.Now()

	err := ea.client.Insert(ctx, accModel)
	if err != nil {
		errMessage := "failed to add account"
		ea.logger.With("address", account.Address).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return accModel.ToEntity(), nil
}

func (ea *ETHAccounts) Update(ctx context.Context, account *entities.ETHAccount) (*entities.ETHAccount, error) {
	accModel := models.NewETHAccount(account)
	accModel.StoreID = ea.storeID
	accModel.UpdatedAt = time.Now()

	err := ea.client.UpdatePK(ctx, accModel)
	if err != nil {
		errMessage := "failed to update account"
		ea.logger.With("address", account.Address).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return accModel.ToEntity(), nil
}

func (ea *ETHAccounts) Delete(ctx context.Context, addr string) error {
	err := ea.client.DeletePK(ctx, &models.ETHAccount{Address: addr, StoreID: ea.storeID})
	if err != nil {
		errMessage := "failed to delete account"
		ea.logger.With("address", addr).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (ea *ETHAccounts) Restore(ctx context.Context, addr string) error {
	accModel := &models.ETHAccount{
		Address: addr,
		StoreID: ea.storeID,
	}
	err := ea.client.UndeletePK(ctx, accModel)
	if err != nil {
		errMessage := "failed to restore account"
		ea.logger.With("address", addr).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (ea *ETHAccounts) Purge(ctx context.Context, addr string) error {
	err := ea.client.ForceDeletePK(ctx, &models.ETHAccount{Address: addr, StoreID: ea.storeID})
	if err != nil {
		errMessage := "failed to permanently delete account"
		ea.logger.With("address", addr).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}
