package aliases

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (i *Aliases) Get(ctx context.Context, registry, aliasKey string) (*entities.Alias, error) {
	return i.db.Get(ctx, registry, aliasKey)
}
