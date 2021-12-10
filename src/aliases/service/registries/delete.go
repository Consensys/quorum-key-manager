package registries

import (
	"context"
)

func (s *Registries) Delete(ctx context.Context, name string) error {
	logger := s.logger.With("name", name)

	err := s.db.Delete(ctx, name)
	if err != nil {
		return err
	}

	logger.Info("alias registry deleted successfully")
	return nil
}
