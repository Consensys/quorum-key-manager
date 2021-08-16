package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c Connector) Info(ctx context.Context) (*entities.StoreInfo, error) {
	return c.store.Info(ctx)
}
