package imports

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

func ImportKeys(ctx context.Context, db database.Keys, mnf *manifest.Manifest, logger log.Logger) error {
	return nil
}
