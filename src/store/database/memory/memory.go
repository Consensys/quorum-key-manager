package memory

import (
	"context"
	"sync"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/database"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
)

type ETH1AccountsDataAgent struct {
	addrToAccounts        map[string]*entities.ETH1Account
	deletedAddrToAccounts map[string]*entities.ETH1Account
	mux                   sync.RWMutex
	logger                *log.Logger
}

var _ database.ETH1Accounts = &ETH1AccountsDataAgent{}

func New(logger *log.Logger) *ETH1AccountsDataAgent {
	return &ETH1AccountsDataAgent{
		mux:                   sync.RWMutex{},
		addrToAccounts:        make(map[string]*entities.ETH1Account),
		deletedAddrToAccounts: make(map[string]*entities.ETH1Account),
		logger:                logger,
	}
}

func (d *ETH1AccountsDataAgent) Get(_ context.Context, addr string) (*entities.ETH1Account, error) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	account, ok := d.addrToAccounts[addr]
	if !ok {
		return nil, errors.NotFoundError("account %s was not found", addr)
	}

	return account, nil
}

func (d *ETH1AccountsDataAgent) GetDeleted(_ context.Context, addr string) (*entities.ETH1Account, error) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	id, ok := d.deletedAddrToAccounts[addr]
	if !ok {
		return nil, errors.NotFoundError("deleted account %s was not found", addr)
	}

	return id, nil
}

func (d *ETH1AccountsDataAgent) GetAll(_ context.Context) ([]*entities.ETH1Account, error) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	accounts := []*entities.ETH1Account{}

	for _, account := range d.addrToAccounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (d *ETH1AccountsDataAgent) GetAllDeleted(_ context.Context) ([]*entities.ETH1Account, error) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	accounts := []*entities.ETH1Account{}

	for _, account := range d.deletedAddrToAccounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (d *ETH1AccountsDataAgent) Add(_ context.Context, account *entities.ETH1Account) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	d.addrToAccounts[account.Address] = account

	return nil
}

func (d *ETH1AccountsDataAgent) AddDeleted(_ context.Context, account *entities.ETH1Account) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.deletedAddrToAccounts[account.Address] = account

	return nil
}

func (d *ETH1AccountsDataAgent) Remove(_ context.Context, addr string) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	delete(d.addrToAccounts, addr)

	return nil
}

func (d *ETH1AccountsDataAgent) RemoveDeleted(_ context.Context, addr string) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	delete(d.deletedAddrToAccounts, addr)

	return nil
}
