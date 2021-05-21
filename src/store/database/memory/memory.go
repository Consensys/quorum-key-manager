package memory

import (
	"context"
	"sync"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/database"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
)

type Database struct {
	addrToID        map[string]string
	deletedAddrToID map[string]string
	mux             sync.RWMutex
	logger          *log.Logger
}

var _ database.Database = &Database{}

func New(logger *log.Logger) *Database {
	return &Database{
		mux:             sync.RWMutex{},
		addrToID:        make(map[string]string),
		deletedAddrToID: make(map[string]string),
		logger:          logger,
	}
}

func (d *Database) GetID(_ context.Context, addr string) (string, error) {
	id, ok := d.addrToID[addr]
	if !ok {
		return "", errors.NotFoundError("account %s was not found", addr)
	}

	return id, nil
}

func (d *Database) GetDeletedID(_ context.Context, addr string) (string, error) {
	id, ok := d.deletedAddrToID[addr]
	if !ok {
		return "", errors.NotFoundError("deleted account %s was not found", addr)
	}

	return id, nil
}

func (d *Database) GetAll(_ context.Context) ([]string, error) {
	addresses := make([]string, len(d.addrToID))

	for address := range d.addrToID {
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (d *Database) GetAllIDs(_ context.Context) ([]string, error) {
	ids := make([]string, len(d.addrToID))

	for _, id := range d.addrToID {
		ids = append(ids, id)
	}

	return ids, nil
}

func (d *Database) GetAllDeleted(_ context.Context) ([]string, error) {
	addresses := make([]string, len(d.deletedAddrToID))

	for address := range d.deletedAddrToID {
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (d *Database) AddID(_ context.Context, addr, id string) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.addrToID[addr] = id

	return nil
}

func (d *Database) AddDeletedID(_ context.Context, addr, id string) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.deletedAddrToID[addr] = id

	return nil
}

func (d *Database) RemoveID(_ context.Context, addr string) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	delete(d.addrToID, addr)

	return nil
}

func (d *Database) RemoveDeletedID(_ context.Context, addr string) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	delete(d.deletedAddrToID, addr)

	return nil
}
