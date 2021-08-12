package models

import (
	"time"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

type Secret struct {
	tableName struct{} `pg:"secrets"` // nolint:unused,structcheck // reason

	ID        string `pg:",pk"`
	Version   string `pg:",pk"`
	StoreID   string `pg:",pk"`
	Tags      map[string]string
	Disabled  bool
	CreatedAt time.Time `pg:"default:now()"`
	UpdatedAt time.Time `pg:"default:now()"`
	DeletedAt time.Time `pg:",soft_delete"`
}

func NewSecret(secret *entities.Secret) *Secret {
	return &Secret{
		ID:        secret.ID,
		Version:   secret.Metadata.Version,
		Tags:      secret.Tags,
		Disabled:  secret.Metadata.Disabled,
		CreatedAt: secret.Metadata.CreatedAt,
		UpdatedAt: secret.Metadata.UpdatedAt,
		DeletedAt: secret.Metadata.DeletedAt,
	}
}

func (s *Secret) ToEntity() *entities.Secret {
	return &entities.Secret{
		ID:   s.ID,
		Tags: s.Tags,
		Metadata: &entities.Metadata{
			Version:   s.Version,
			Disabled:  s.Disabled,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
			DeletedAt: s.DeletedAt,
		},
	}
}
