package memory

import (
	"context"
	"sync"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/database"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
)

type ETH1Accounts struct {
	addrToAccounts        map[string]*entities.ETH1Account
	deletedAddrToAccounts map[string]*entities.ETH1Account
	mux                   sync.RWMutex
	logger                *log.Logger
}

var _ database.ETH1Accounts = &ETH1Accounts{}

func New(logger *log.Logger) *ETH1Accounts {
	return &ETH1Accounts{
		mux:                   sync.RWMutex{},
		addrToAccounts:        make(map[string]*entities.ETH1Account),
		deletedAddrToAccounts: make(map[string]*entities.ETH1Account),
		logger:                logger,
	}
}

func (d *ETH1Accounts) GetAccount(_ context.Context, addr string) (*entities.ETH1Account, error) {
	account, ok := d.addrToAccounts[addr]
	if !ok {
		return nil, errors.NotFoundError("account %s was not found", addr)
	}

	return account, nil
}

func (d *ETH1Accounts) GetDeletedAccount(_ context.Context, addr string) (*entities.ETH1Account, error) {
	id, ok := d.deletedAddrToAccounts[addr]
	if !ok {
		return nil, errors.NotFoundError("deleted account %s was not found", addr)
	}

	return id, nil
}

func (d *ETH1Accounts) GetAllAccounts(_ context.Context) ([]*entities.ETH1Account, error) {
	accounts := []*entities.ETH1Account{}

	for _, account := range d.addrToAccounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (d *ETH1Accounts) GetAllDeletedAccounts(_ context.Context) ([]*entities.ETH1Account, error) {
	accounts := []*entities.ETH1Account{}

	for _, account := range d.deletedAddrToAccounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (d *ETH1Accounts) AddAccount(_ context.Context, account *entities.ETH1Account) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.addrToAccounts[account.Address] = account

	return nil
}

func (d *ETH1Accounts) AddDeletedAccount(_ context.Context, account *entities.ETH1Account) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.deletedAddrToAccounts[account.Address] = account

	return nil
}

func (d *ETH1Accounts) RemoveAccount(_ context.Context, addr string) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	delete(d.addrToAccounts, addr)

	return nil
}

func (d *ETH1Accounts) RemoveDeletedAccount(_ context.Context, addr string) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	delete(d.deletedAddrToAccounts, addr)

	return nil
}
