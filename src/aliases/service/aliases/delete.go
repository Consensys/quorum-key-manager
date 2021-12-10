package aliases

import "context"

func (s *Aliases) Delete(ctx context.Context, registry, key string) error {
	logger := s.logger.With("registry", registry, "key", key)

	err := s.db.Delete(ctx, registry, key)
	if err != nil {
		return err
	}

	logger.Info("alias deleted successfully")
	return nil
}
