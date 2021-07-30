package postgres

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/stores/store/database/models"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"time"

	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
	"github.com/consensys/quorum-key-manager/src/stores/store/database"

	"github.com/consensys/quorum-key-manager/pkg/errors"
)

type ETH1Accounts struct {
	logger log.Logger
	db     postgres.Client
}

var _ database.ETH1Accounts = &ETH1Accounts{}

func NewETH1Accounts(logger log.Logger, db postgres.Client) *ETH1Accounts {
	return &ETH1Accounts{
		logger: logger,
		db:     db,
	}
}

func (d *ETH1Accounts) Get(ctx context.Context, addr string) (*entities.ETH1Account, error) {
	eth1Acc := &models.ETH1Account{Address: addr}

	err := d.db.SelectPK(ctx, eth1Acc)
	if err != nil {
		errMessage := "failed to get account"
		d.logger.With("address", addr).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return eth1Acc.ToEntity(), nil
}

func (d *ETH1Accounts) GetDeleted(ctx context.Context, addr string) (*entities.ETH1Account, error) {
	eth1Acc := &models.ETH1Account{Address: addr}

	err := d.db.SelectDeletedPK(ctx, eth1Acc)
	if err != nil {
		errMessage := "failed to get deleted account"
		d.logger.With("address", addr).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return eth1Acc.ToEntity(), nil
}

func (d *ETH1Accounts) GetAll(ctx context.Context) ([]*entities.ETH1Account, error) {
	var eth1Accs []*models.ETH1Account

	err := d.db.Select(ctx, &eth1Accs)
	if err != nil {
		errMessage := "failed to get all accounts"
		d.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	var accounts []*entities.ETH1Account
	for _, acc := range eth1Accs {
		accounts = append(accounts, acc.ToEntity())
	}

	return accounts, nil
}

func (d *ETH1Accounts) GetAllDeleted(ctx context.Context) ([]*entities.ETH1Account, error) {
	var eth1Accs []*models.ETH1Account

	err := d.db.SelectDeleted(ctx, &eth1Accs)
	if err != nil {
		errMessage := "failed to get all deleted accounts"
		d.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	var accounts []*entities.ETH1Account
	for _, acc := range eth1Accs {
		accounts = append(accounts, acc.ToEntity())
	}

	return accounts, nil
}

func (d *ETH1Accounts) Add(ctx context.Context, account *entities.ETH1Account) (*entities.ETH1Account, error) {
	accModel := models.NewETH1Account(account)

	err := d.db.Insert(ctx, accModel)
	if err != nil {
		errMessage := "failed to add account"
		d.logger.With("address", account.Address).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return accModel.ToEntity(), nil
}

func (d *ETH1Accounts) Update(ctx context.Context, account *entities.ETH1Account) (*entities.ETH1Account, error) {
	accModel := models.NewETH1Account(account)

	err := d.db.UpdatePK(ctx, accModel)
	if err != nil {
		errMessage := "failed to update account"
		d.logger.With("address", account.Address).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return accModel.ToEntity(), nil
}

func (d *ETH1Accounts) Delete(ctx context.Context, addr string) error {
	err := d.db.DeletePK(ctx, &models.ETH1Account{Address: addr})
	if err != nil {
		errMessage := "failed to delete account"
		d.logger.With("address", addr).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (d *ETH1Accounts) Restore(ctx context.Context, account *entities.ETH1Account) error {
	account.Metadata.DeletedAt = time.Time{}
	_, err := d.Update(ctx, account)
	if err != nil {
		errMessage := "failed to restore account"
		d.logger.With("address", account.Address).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (d *ETH1Accounts) Purge(ctx context.Context, addr string) error {
	err := d.db.ForceDeletePK(ctx, &models.ETH1Account{Address: addr})
	if err != nil {
		errMessage := "failed to permanently delete account"
		d.logger.With("address", addr).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}
