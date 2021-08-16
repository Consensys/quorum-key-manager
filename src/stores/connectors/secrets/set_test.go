package secrets

import (
	"context"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	mock2 "github.com/consensys/quorum-key-manager/src/stores/database/mock"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSetSecret(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockSecretStore(ctrl)
	db := mock2.NewMockSecrets(ctrl)
	logger := testutils.NewMockLogger(ctrl)

	connector := NewConnector(store, db, nil, logger)

	t.Run("should set secret successfully", func(t *testing.T) {
		secret := testutils2.FakeSecret()
		attributes := &entities.Attributes{}

		store.EXPECT().Set(gomock.Any(), secret.ID, secret.Value, attributes).Return(secret, nil)

		db.EXPECT().Add(gomock.Any(), secret).Return(secret, nil)

		rSecret, err := connector.Set(ctx, secret.ID, secret.Value, attributes)

		assert.NoError(t, err)
		assert.Equal(t, rSecret, secret)
	})

	t.Run("should fail to delete secret if store fail to set", func(t *testing.T) {
		secret := testutils2.FakeSecret()
		attributes := &entities.Attributes{}
		expectedErr := errors.UnauthorizedError("not authorized")

		store.EXPECT().Set(gomock.Any(), secret.ID, secret.Value, attributes).Return(nil, expectedErr)

		_, err := connector.Set(ctx, secret.ID, secret.Value, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to set secret if db fail to add", func(t *testing.T) {
		secret := testutils2.FakeSecret()
		attributes := &entities.Attributes{}
		expectedErr := errors.NotFoundError("not found")

		store.EXPECT().Set(gomock.Any(), secret.ID, secret.Value, attributes).Return(secret, nil)

		db.EXPECT().Add(gomock.Any(), secret).Return(nil, expectedErr)

		_, err := connector.Set(ctx, secret.ID, secret.Value, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
