package keys

import (
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

type Connector struct {
	store        stores.KeyStore
	db           database.Keys
	logger       log.Logger
	authorizator auth.Authorizator
}

var _ stores.KeyStore = Connector{}

func NewConnector(store stores.KeyStore, db database.Keys, authorizator auth.Authorizator, logger log.Logger) *Connector {
	return &Connector{
		store:        store,
		db:           db,
		logger:       logger,
		authorizator: authorizator,
	}
}

func isSupportedAlgo(alg *entities.Algorithm) bool {
	if alg.Type == entities.Ecdsa && alg.EllipticCurve == entities.Secp256k1 {
		return true
	}

	if alg.Type == entities.Eddsa && alg.EllipticCurve == entities.Babyjubjub {
		return true
	}

	return false
}
