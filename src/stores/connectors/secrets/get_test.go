package secrets

import (
	"context"
	"fmt"
	"testing"

	mock3 "github.com/consensys/quorum-key-manager/src/auth/mock"
	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	mock2 "github.com/consensys/quorum-key-manager/src/stores/database/mock"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetSecret(t *testing.T) {
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

	t.Run("should get secret successfully", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(secret, nil)
		store.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(secret, nil)

		rSecret, err := connector.Get(ctx, secret.ID, secret.Metadata.Version)

		assert.NoError(t, err)
		assert.Equal(t, secret, rSecret)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret}).Return(expectedErr)

		_, err := connector.Get(ctx, secret.ID, secret.Metadata.Version)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to get secret if db fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil, expectedErr)

		_, err := connector.Get(ctx, secret.ID, secret.Metadata.Version)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to get secret value", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(secret, nil)
		store.EXPECT().Get(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil, expectedErr)

		_, err := connector.Get(ctx, secret.ID, secret.Metadata.Version)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}

func TestGetDeletedSecret(t *testing.T) {
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

	t.Run("should get deleted secret successfully", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), secret.ID, secret.Metadata.Version).Return(secret, nil)

		rSecret, err := connector.GetDeleted(ctx, secret.ID, secret.Metadata.Version)

		assert.NoError(t, err)
		assert.Equal(t, secret, rSecret)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret}).Return(expectedErr)

		_, err := connector.GetDeleted(ctx, secret.ID, secret.Metadata.Version)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to get deleted secret if db fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), secret.ID, secret.Metadata.Version).Return(nil, expectedErr)

		_, err := connector.GetDeleted(ctx, secret.ID, secret.Metadata.Version)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
