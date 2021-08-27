package secrets

import (
	"context"
	"fmt"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	mock3 "github.com/consensys/quorum-key-manager/src/auth/mock"
	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	mock2 "github.com/consensys/quorum-key-manager/src/stores/database/mock"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSetSecret(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	secret := testutils2.FakeSecret()
	attributes := testutils2.FakeAttributes()
	expectedErr := fmt.Errorf("error")

	store := mock.NewMockSecretStore(ctrl)
	db := mock2.NewMockSecrets(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	t.Run("should set secret successfully", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceSecret}).Return(nil)
		store.EXPECT().Set(gomock.Any(), secret.ID, secret.Value, attributes).Return(secret, nil)
		db.EXPECT().Add(gomock.Any(), secret).Return(secret, nil)

		rSecret, err := connector.Set(ctx, secret.ID, secret.Value, attributes)

		assert.NoError(t, err)
		assert.Equal(t, rSecret, secret)
	})

	t.Run("should create key successfully if it already exists in the vault", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceSecret}).Return(nil)
		store.EXPECT().Set(gomock.Any(), secret.ID, secret.Value, attributes).Return(nil, errors.AlreadyExistsError("error"))
		store.EXPECT().Get(gomock.Any(), secret.ID, "").Return(secret, nil)
		db.EXPECT().Add(gomock.Any(), secret).Return(secret, nil)

		rSecret, err := connector.Set(ctx, secret.ID, secret.Value, attributes)

		assert.NoError(t, err)
		assert.Equal(t, rSecret, secret)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceSecret}).Return(expectedErr)

		_, err := connector.Set(ctx, secret.ID, secret.Value, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to delete secret if store fail to set", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceSecret}).Return(nil)
		store.EXPECT().Set(gomock.Any(), secret.ID, secret.Value, attributes).Return(nil, expectedErr)

		_, err := connector.Set(ctx, secret.ID, secret.Value, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to set secret if db fail to add", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceSecret}).Return(nil)
		store.EXPECT().Set(gomock.Any(), secret.ID, secret.Value, attributes).Return(secret, nil)
		db.EXPECT().Add(gomock.Any(), secret).Return(nil, expectedErr)

		_, err := connector.Set(ctx, secret.ID, secret.Value, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
