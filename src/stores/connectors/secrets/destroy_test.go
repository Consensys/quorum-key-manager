package secrets

import (
	"context"
	"fmt"
	"testing"

	mock3 "github.com/consensys/quorum-key-manager/src/auth/mock"
	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	mock2 "github.com/consensys/quorum-key-manager/src/stores/database/mock"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDestroySecret(t *testing.T) {
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

	t.Run("should destroy secret successfully", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDestroy, Resource: types.ResourceSecret}).Return(nil)
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil, errors.NotFoundError("error"))
		db.EXPECT().GetDeleted(gomock.Any(), secret.ID, secret.Metadata.Version).Return(secret, nil)
		db.EXPECT().Purge(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil)
		store.EXPECT().Destroy(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil)

		err := connector.Destroy(ctx, secret.ID, secret.Metadata.Version)

		assert.NoError(t, err)
	})

	t.Run("should be idempotent if secret is not found and deleted", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDestroy, Resource: types.ResourceSecret}).Return(nil)
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil, errors.NotFoundError("error"))

		db.EXPECT().GetDeleted(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil, expectedErr)

		err := connector.Destroy(ctx, secret.ID, secret.Metadata.Version)

		assert.Nil(t, err)
	})

	t.Run("should destroy secret successfully, ignoring not supported error", func(t *testing.T) {
		rErr := errors.NotSupportedError("not supported")

		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDestroy, Resource: types.ResourceSecret}).Return(nil)
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil, errors.NotFoundError("error"))
		db.EXPECT().GetDeleted(gomock.Any(), secret.ID, secret.Metadata.Version).Return(secret, nil)
		db.EXPECT().Purge(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil)
		store.EXPECT().Destroy(gomock.Any(), secret.ID, secret.Metadata.Version).Return(rErr)

		err := connector.Destroy(ctx, secret.ID, secret.Metadata.Version)

		assert.NoError(t, err)
	})

	t.Run("should fail to destroy secret with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDestroy, Resource: types.ResourceSecret}).Return(expectedErr)

		err := connector.Destroy(ctx, secret.ID, secret.Metadata.Version)

		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to destroy secret if db fails to purge", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDestroy, Resource: types.ResourceSecret}).Return(nil)
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil, errors.NotFoundError("error"))

		db.EXPECT().GetDeleted(gomock.Any(), secret.ID, secret.Metadata.Version).Return(secret, nil)
		db.EXPECT().Purge(gomock.Any(), secret.ID, secret.Metadata.Version).Return(expectedErr)

		err := connector.Destroy(ctx, secret.ID, secret.Metadata.Version)

		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to destroy secret if store fails to destroy", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDestroy, Resource: types.ResourceSecret}).Return(nil)
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil, errors.NotFoundError("error"))
		db.EXPECT().GetDeleted(gomock.Any(), secret.ID, secret.Metadata.Version).Return(secret, nil)
		db.EXPECT().Purge(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil)
		store.EXPECT().Destroy(gomock.Any(), secret.ID, secret.Metadata.Version).Return(expectedErr)

		err := connector.Destroy(ctx, secret.ID, secret.Metadata.Version)

		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to destroy secret if secret already exists", func(t *testing.T) {
		expectedErr = errors.StatusConflictError("secret %s must be deleted first", secret.ID)
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDestroy, Resource: types.ResourceSecret}).Return(nil)
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(secret, nil)
		store.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(secret, nil)

		err := connector.Destroy(ctx, secret.ID, secret.Metadata.Version)

		assert.Equal(t, err, expectedErr)
	})
}
