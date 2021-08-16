package eth1

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c Connector) Info(_ context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}
