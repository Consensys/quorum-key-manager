package keys

import (
	"context"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	mock2 "github.com/consensys/quorum-key-manager/src/stores/database/mock"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetKey(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockKeys(ctrl)
	logger := testutils.NewMockLogger(ctrl)

	connector := NewConnector(store, db, nil, logger)

	t.Run("should get key successfully", func(t *testing.T) {
		key := testutils2.FakeKey()

		db.EXPECT().Get(gomock.Any(), key.ID).Return(key, nil)

		rKey, err := connector.Get(ctx, key.ID)

		assert.NoError(t, err)
		assert.Equal(t, key, rKey)
	})

	t.Run("should fail to get key if db fails", func(t *testing.T) {
		key := testutils2.FakeKey()
		expectedErr := errors.PostgresError("cannot connect")

		db.EXPECT().Get(gomock.Any(), key.ID).Return(nil, expectedErr)

		_, err := connector.Get(ctx, key.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}

func TestGetDeletedKey(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockKeys(ctrl)
	logger := testutils.NewMockLogger(ctrl)

	connector := NewConnector(store, db, nil, logger)

	t.Run("should get deleted key successfully", func(t *testing.T) {
		key := testutils2.FakeKey()

		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(key, nil)

		rKey, err := connector.GetDeleted(ctx, key.ID)

		assert.NoError(t, err)
		assert.Equal(t, key, rKey)
	})

	t.Run("should fail to get deleted key if db fails", func(t *testing.T) {
		key := testutils2.FakeKey()
		expectedErr := errors.PostgresError("cannot connect")

		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(nil, expectedErr)

		_, err := connector.GetDeleted(ctx, key.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
