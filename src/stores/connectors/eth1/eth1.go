package eth1

import (
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

type Connector struct {
	store        stores.KeyStore
	logger       log.Logger
	db           database.ETH1Accounts
	authorizator auth.Authorizator
}

var _ stores.Eth1Store = Connector{}

var eth1Algo = &entities.Algorithm{
	Type:          entities.Ecdsa,
	EllipticCurve: entities.Secp256k1,
}

func NewConnector(store stores.KeyStore, db database.ETH1Accounts, authorizator auth.Authorizator, logger log.Logger) *Connector {
	return &Connector{
		store:        store,
		logger:       logger,
		db:           db,
		authorizator: authorizator,
	}
}
