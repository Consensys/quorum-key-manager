package aliases

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
)

func (s *Aliases) Delete(ctx context.Context, registry, key string, userInfo *auth.UserInfo) error {
	logger := s.logger.With("registry", registry, "key", key)

	err := s.db.Delete(ctx, registry, key, userInfo.Tenant)
	if err != nil {
		errMessage := "failed to delete alias"
		logger.WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	logger.Info("alias deleted successfully")
	return nil
}
