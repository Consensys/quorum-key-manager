package eth1

import (
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/ethereum/go-ethereum/crypto"
)

type Connector struct {
	store  stores.KeyStore
	logger log.Logger
	db     database.ETH1Accounts
}

var _ stores.Eth1Store = Connector{}

var eth1Algo = &entities.Algorithm{
	Type:          entities.Ecdsa,
	EllipticCurve: entities.Secp256k1,
}

func NewConnector(store stores.KeyStore, db database.ETH1Accounts, logger log.Logger) *Connector {
	return &Connector{
		store:  store,
		logger: logger,
		db:     db,
	}
}

func newEth1Account(key *entities.Key, attr *entities.Attributes) *entities.ETH1Account {
	pubKey, _ := crypto.UnmarshalPubkey(key.PublicKey)
	return &entities.ETH1Account{
		KeyID:               key.ID,
		Address:             crypto.PubkeyToAddress(*pubKey),
		Tags:                attr.Tags,
		PublicKey:           key.PublicKey,
		CompressedPublicKey: crypto.CompressPubkey(pubKey),
		Metadata: &entities.Metadata{
			Disabled:  key.Metadata.Disabled,
			CreatedAt: key.Metadata.CreatedAt,
			UpdatedAt: key.Metadata.UpdatedAt,
		},
	}
}
