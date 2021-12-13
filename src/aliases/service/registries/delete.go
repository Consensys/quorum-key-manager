package registries

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
)

func (s *Registries) Delete(ctx context.Context, name string, userInfo *auth.UserInfo) error {
	logger := s.logger.With("name", name)

	err := s.db.Delete(ctx, name, userInfo.Tenant)
	if err != nil {
		errMessage := "failed to delete registry"
		logger.WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	logger.Info("alias registry deleted successfully")
	return nil
}
