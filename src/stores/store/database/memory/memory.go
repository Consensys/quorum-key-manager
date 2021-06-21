package memory

import (
	"context"
	"sync"

	"github.com/consensysquorum/quorum-key-manager/pkg/log"

	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"

	"github.com/consensysquorum/quorum-key-manager/src/stores/store/database"

	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
)

type ETH1Accounts struct {
	addrToAccounts        map[string]*entities.ETH1Account
	deletedAddrToAccounts map[string]*entities.ETH1Account
	mux                   sync.RWMutex
	logger                log.Logger
}

var _ database.ETH1Accounts = &ETH1Accounts{}

func New(logger log.Logger) *ETH1Accounts {
	return &ETH1Accounts{
		mux:                   sync.RWMutex{},
		addrToAccounts:        make(map[string]*entities.ETH1Account),
		deletedAddrToAccounts: make(map[string]*entities.ETH1Account),
		logger:                logger,
	}
}

func (d *ETH1Accounts) Get(_ context.Context, addr string) (*entities.ETH1Account, error) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	account, ok := d.addrToAccounts[addr]
	if !ok {
		errMessage := "account was not found"
		d.logger.Error(errMessage, "account", addr)
		return nil, errors.NotFoundError(errMessage)
	}

	return account, nil
}

func (d *ETH1Accounts) GetDeleted(_ context.Context, addr string) (*entities.ETH1Account, error) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	id, ok := d.deletedAddrToAccounts[addr]
	if !ok {
		errMessage := "deleted account was not found"
		d.logger.Error(errMessage, "account", addr)
		return nil, errors.NotFoundError(errMessage)
	}

	return id, nil
}

func (d *ETH1Accounts) GetAll(_ context.Context) ([]*entities.ETH1Account, error) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	accounts := []*entities.ETH1Account{}

	for _, account := range d.addrToAccounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (d *ETH1Accounts) GetAllDeleted(_ context.Context) ([]*entities.ETH1Account, error) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	accounts := []*entities.ETH1Account{}

	for _, account := range d.deletedAddrToAccounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (d *ETH1Accounts) Add(_ context.Context, account *entities.ETH1Account) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	if _, ok := d.addrToAccounts[account.Address.Hex()]; ok {
		errMessage := "account already exists"
		d.logger.Error(errMessage, "account", account.Address)
		return errors.AlreadyExistsError(errMessage)
	}

	if _, ok := d.deletedAddrToAccounts[account.Address.Hex()]; ok {
		errMessage := "account is currently deleted. Please restore it instead"
		d.logger.Error(errMessage, "account", account.Address)
		return errors.AlreadyExistsError(errMessage)
	}

	d.addrToAccounts[account.Address.Hex()] = account

	return nil
}

func (d *ETH1Accounts) Update(_ context.Context, account *entities.ETH1Account) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	if _, ok := d.addrToAccounts[account.Address.Hex()]; !ok {
		errMessage := "account does not exists"
		d.logger.Error(errMessage, "account", account.Address)
		return errors.NotFoundError(errMessage)
	}

	d.addrToAccounts[account.Address.Hex()] = account

	return nil
}

func (d *ETH1Accounts) AddDeleted(_ context.Context, account *entities.ETH1Account) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	d.deletedAddrToAccounts[account.Address.Hex()] = account

	return nil
}

func (d *ETH1Accounts) Remove(_ context.Context, addr string) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	delete(d.addrToAccounts, addr)

	return nil
}

func (d *ETH1Accounts) RemoveDeleted(_ context.Context, addr string) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	delete(d.deletedAddrToAccounts, addr)

	return nil
}
