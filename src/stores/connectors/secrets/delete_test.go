package secrets

import (
	"context"
	"testing"

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

	store := mock.NewMockSecretStore(ctrl)
	db := mock2.NewMockSecrets(ctrl)
	logger := testutils.NewMockLogger(ctrl)

	connector := NewConnector(store, db, logger)

	db.EXPECT().RunInTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, persist func(dbtx database.Secrets) error) error {
			return persist(db)
		}).AnyTimes()

	t.Run("should delete secret successfully", func(t *testing.T) {
		secret := testutils2.FakeSecret()

		db.EXPECT().Delete(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil)

		store.EXPECT().Delete(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil)

		err := connector.Delete(ctx, secret.ID, secret.Metadata.Version)

		assert.NoError(t, err)
	})

	t.Run("should delete secret successfully, ignoring not supported error", func(t *testing.T) {
		secret := testutils2.FakeSecret()
		rErr := errors.NotSupportedError("not supported")

		db.EXPECT().Delete(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil)

		store.EXPECT().Delete(gomock.Any(), secret.ID, secret.Metadata.Version).Return(rErr)

		err := connector.Delete(ctx, secret.ID, secret.Metadata.Version)

		assert.NoError(t, err)
	})

	t.Run("should fail to delete secret if db fail to delete", func(t *testing.T) {
		secret := testutils2.FakeSecret()
		expectedErr := errors.NotFoundError("not found")

		db.EXPECT().Delete(gomock.Any(), secret.ID, secret.Metadata.Version).Return(expectedErr)

		err := connector.Delete(ctx, secret.ID, secret.Metadata.Version)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to delete secret if store fail to delete", func(t *testing.T) {
		secret := testutils2.FakeSecret()
		expectedErr := errors.UnauthorizedError("not authorized")

		db.EXPECT().Delete(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil)

		store.EXPECT().Delete(gomock.Any(), secret.ID, secret.Metadata.Version).Return(expectedErr)

		err := connector.Delete(ctx, secret.ID, secret.Metadata.Version)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
