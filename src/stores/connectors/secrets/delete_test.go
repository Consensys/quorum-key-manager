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

func TestDeleteSecret(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	t.Run("should delete secret successfully", func(t *testing.T) {
		secret := testutils2.FakeSecret()

		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDelete, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().Delete(gomock.Any(), secret.ID).Return(nil)
		store.EXPECT().Delete(gomock.Any(), secret.ID).Return(nil)

		err := connector.Delete(ctx, secret.ID)

		assert.NoError(t, err)
	})

	t.Run("should delete secret successfully, ignoring not supported error", func(t *testing.T) {
		secret := testutils2.FakeSecret()
		rErr := errors.NotSupportedError("not supported")

		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDelete, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().Delete(gomock.Any(), secret.ID).Return(nil)
		store.EXPECT().Delete(gomock.Any(), secret.ID).Return(rErr)

		err := connector.Delete(ctx, secret.ID)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		secret := testutils2.FakeSecret()

		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDelete, Resource: types.ResourceSecret}).Return(expectedErr)

		err := connector.Delete(ctx, secret.ID)

		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to delete secret if db fail to delete", func(t *testing.T) {
		secret := testutils2.FakeSecret()

		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDelete, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().Delete(gomock.Any(), secret.ID).Return(expectedErr)

		err := connector.Delete(ctx, secret.ID)

		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to delete secret if store fail to delete", func(t *testing.T) {
		secret := testutils2.FakeSecret()

		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionDelete, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().Delete(gomock.Any(), secret.ID).Return(nil)
		store.EXPECT().Delete(gomock.Any(), secret.ID).Return(expectedErr)

		err := connector.Delete(ctx, secret.ID)

		assert.Equal(t, err, expectedErr)
	})
}
