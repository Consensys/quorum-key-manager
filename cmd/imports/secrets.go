package imports

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

func ImportSecrets(ctx context.Context, storeName string, db database.Database, vault stores.SecretStore) error {

	return nil
}
