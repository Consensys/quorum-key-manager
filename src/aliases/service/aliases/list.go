package aliases

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (i *Aliases) List(ctx context.Context, registry string) ([]entities.Alias, error) {
	return i.db.List(ctx, registry)
}
