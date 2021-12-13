package registries

import (
	"context"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
)

func (s *Registries) Delete(ctx context.Context, name string, userInfo *auth.UserInfo) error {
	logger := s.logger.With("name", name)

	err := s.db.Delete(ctx, name, userInfo.Tenant)
	if err != nil {
		return err
	}

	logger.Info("alias registry deleted successfully")
	return nil
}
