package secrets

import (
	"context"
	"fmt"
	"testing"

	"github.com/consensys/quorum-key-manager/src/auth/entities"
	mock3 "github.com/consensys/quorum-key-manager/src/auth/mock"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	mock2 "github.com/consensys/quorum-key-manager/src/stores/database/mock"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRestoreSecret(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	secret := testutils2.FakeSecret()
	expectedErr := fmt.Errorf("error")

	store := mock.NewMockSecretStore(ctrl)
	db := mock2.NewMockSecrets(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	db.EXPECT().RunInTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, persist func(dbtx database.Secrets) error) error {
			return persist(db)
		}).AnyTimes()

	t.Run("should restore secret successfully", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceSecret}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceSecret}).Return(nil)
		db.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil, errors.NotFoundError("error"))
		db.EXPECT().GetLatestVersion(gomock.Any(), secret.ID, false).Return(secret.Metadata.Version, nil)
		db.EXPECT().GetDeleted(gomock.Any(), secret.ID).Return(secret, nil)
		db.EXPECT().Restore(gomock.Any(), secret.ID).Return(nil)
		store.EXPECT().Restore(gomock.Any(), secret.ID).Return(nil)

		err := connector.Restore(ctx, secret.ID)

		assert.NoError(t, err)
	})

	t.Run("should be idempotent if secret exists", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceSecret}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceSecret}).Return(nil)
		db.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(secret, nil)
		db.EXPECT().GetLatestVersion(gomock.Any(), secret.ID, false).Return(secret.Metadata.Version, nil)
		store.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(secret, nil)

		err := connector.Restore(ctx, secret.ID)

		assert.NoError(t, err)
	})

	t.Run("should restore secret successfully, ignoring not supported error", func(t *testing.T) {
		rErr := errors.NotSupportedError("not supported")

		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceSecret}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceSecret}).Return(nil)
		db.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil, errors.NotFoundError("error"))
		db.EXPECT().GetLatestVersion(gomock.Any(), secret.ID, false).Return(secret.Metadata.Version, nil)
		db.EXPECT().GetDeleted(gomock.Any(), secret.ID).Return(secret, nil)
		db.EXPECT().Restore(gomock.Any(), secret.ID).Return(nil)
		store.EXPECT().Restore(gomock.Any(), secret.ID).Return(rErr)

		err := connector.Restore(ctx, secret.ID)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceSecret}).Return(expectedErr)

		err := connector.Restore(ctx, secret.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to restore secret if secret is not found and not deleted", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceSecret}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceSecret}).Return(nil)
		db.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil, errors.NotFoundError("error"))
		db.EXPECT().GetLatestVersion(gomock.Any(), secret.ID, false).Return(secret.Metadata.Version, nil)
		db.EXPECT().GetDeleted(gomock.Any(), secret.ID).Return(secret, expectedErr)

		err := connector.Restore(ctx, secret.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to restore secret if db fail to restore", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceSecret}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceSecret}).Return(nil)
		db.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil, errors.NotFoundError("error"))
		db.EXPECT().GetLatestVersion(gomock.Any(), secret.ID, false).Return(secret.Metadata.Version, nil)
		db.EXPECT().GetDeleted(gomock.Any(), secret.ID).Return(secret, nil)
		db.EXPECT().Restore(gomock.Any(), secret.ID).Return(expectedErr)

		err := connector.Restore(ctx, secret.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to restore secret if store fail to restore", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceSecret}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceSecret}).Return(nil)
		db.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil, errors.NotFoundError("error"))
		db.EXPECT().GetLatestVersion(gomock.Any(), secret.ID, false).Return(secret.Metadata.Version, nil)
		db.EXPECT().GetDeleted(gomock.Any(), secret.ID).Return(secret, nil)
		db.EXPECT().Restore(gomock.Any(), secret.ID).Return(nil)
		store.EXPECT().Restore(gomock.Any(), secret.ID).Return(expectedErr)

		err := connector.Restore(ctx, secret.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
