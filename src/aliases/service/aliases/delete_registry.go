package aliases

import "context"

func (i *Aliases) DeleteRegistry(ctx context.Context, registry string) error {
	logger := i.logger.With("registry", registry)

	err := i.db.DeleteRegistry(ctx, registry)
	if err != nil {
		return err
	}

	logger.Info("registry deleted successfully")
	return nil
}
