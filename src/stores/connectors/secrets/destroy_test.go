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
		db.EXPECT().GetDeleted(gomock.Any(), secret.ID).Return(secret, nil)
		db.EXPECT().Purge(gomock.Any(), secret.ID).Return(nil)
		store.EXPECT().Destroy(gomock.Any(), secret.ID).Return(nil)

		err := connector.Destroy(ctx, secret.ID)

		assert.NoError(t, err)
	})

	t.Run("should destroy secret successfully, ignoring not supported error", func(t *testing.T) {
		rErr := errors.NotSupportedError("not supported")

		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDestroy, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), secret.ID).Return(secret, nil)
		db.EXPECT().Purge(gomock.Any(), secret.ID).Return(nil)
		store.EXPECT().Destroy(gomock.Any(), secret.ID).Return(rErr)

		err := connector.Destroy(ctx, secret.ID)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDestroy, Resource: types.ResourceSecret}).Return(expectedErr)

		err := connector.Destroy(ctx, secret.ID)

		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to destroy secret if secret is not deleted", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDestroy, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), secret.ID).Return(secret, expectedErr)

		err := connector.Destroy(ctx, secret.ID)

		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to destroy secret if db fail to purge", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDestroy, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), secret.ID).Return(secret, nil)
		db.EXPECT().Purge(gomock.Any(), secret.ID).Return(expectedErr)

		err := connector.Destroy(ctx, secret.ID)

		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to destroy secret if store fail to destroy", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDestroy, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), secret.ID).Return(secret, nil)
		db.EXPECT().Purge(gomock.Any(), secret.ID).Return(nil)
		store.EXPECT().Destroy(gomock.Any(), secret.ID).Return(expectedErr)

		err := connector.Destroy(ctx, secret.ID)

		assert.Equal(t, err, expectedErr)
	})
}
