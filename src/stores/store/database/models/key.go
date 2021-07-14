package models

import (
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"time"
)

type Key struct {
	tableName struct{} `pg:"keys"`

	ID               string
	PublicKey        []byte
	SigningAlgorithm string
	EllipticCurve    string
	Tags             map[string]string
	Annotations      map[string]string
	Disabled         bool
	ExpireAt         time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time `pg:",soft_delete"`
	DestroyedAt      time.Time
}

func NewKey(k *entities.Key) *Key {
	return &Key{
		ID:               k.ID,
		PublicKey:        k.PublicKey,
		SigningAlgorithm: string(k.Algo.Type),
		EllipticCurve:    string(k.Algo.EllipticCurve),
		Tags:             k.Tags,
		Annotations:      k.Annotations,
		Disabled:         k.Metadata.Disabled,
		ExpireAt:         k.Metadata.ExpireAt,
		CreatedAt:        k.Metadata.CreatedAt,
		UpdatedAt:        k.Metadata.UpdatedAt,
		DeletedAt:        k.Metadata.DeletedAt,
		DestroyedAt:      k.Metadata.DestroyedAt,
	}
}

func (k *Key) ToEntity() *entities.Key {
	return &entities.Key{
		ID:        k.ID,
		PublicKey: k.PublicKey,
		Algo: &entities.Algorithm{
			Type:          entities.KeyType(k.SigningAlgorithm),
			EllipticCurve: entities.Curve(k.EllipticCurve),
		},
		Tags:        k.Tags,
		Annotations: k.Annotations,
		Metadata: &entities.Metadata{
			Disabled:    k.Disabled,
			ExpireAt:    k.ExpireAt,
			CreatedAt:   k.CreatedAt,
			UpdatedAt:   k.UpdatedAt,
			DeletedAt:   k.DeletedAt,
			DestroyedAt: k.DestroyedAt,
		},
	}
}
