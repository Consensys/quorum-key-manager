package aliases

import "context"

func (i *Aliases) Delete(ctx context.Context, registry, aliasKey string) error {
	logger := i.logger.With("registry", registry, "key", aliasKey)

	err := i.db.Delete(ctx, registry, aliasKey)
	if err != nil {
		return err
	}

	logger.Info("alias deleted successfully")
	return nil
}
