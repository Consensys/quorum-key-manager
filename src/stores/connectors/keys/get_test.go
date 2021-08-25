package keys

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

func TestGetKey(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := testutils2.FakeKey()
	expectedErr := fmt.Errorf("error")

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockKeys(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	t.Run("should get key successfully", func(t *testing.T) {
		auth.EXPECT().Check(&types.Operation{Action: types.ActionRead, Resource: types.ResourceKey}).Return(nil)
		db.EXPECT().Get(gomock.Any(), key.ID).Return(key, nil)

		rKey, err := connector.Get(ctx, key.ID)

		assert.NoError(t, err)
		assert.Equal(t, key, rKey)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().Check(&types.Operation{Action: types.ActionRead, Resource: types.ResourceKey}).Return(expectedErr)

		_, err := connector.Get(ctx, key.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to get key if db fails", func(t *testing.T) {
		auth.EXPECT().Check(&types.Operation{Action: types.ActionRead, Resource: types.ResourceKey}).Return(nil)
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

	key := testutils2.FakeKey()
	expectedErr := fmt.Errorf("error")

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockKeys(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	t.Run("should get deleted key successfully", func(t *testing.T) {
		auth.EXPECT().Check(&types.Operation{Action: types.ActionRead, Resource: types.ResourceKey}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(key, nil)

		rKey, err := connector.GetDeleted(ctx, key.ID)

		assert.NoError(t, err)
		assert.Equal(t, key, rKey)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().Check(&types.Operation{Action: types.ActionRead, Resource: types.ResourceKey}).Return(expectedErr)

		_, err := connector.GetDeleted(ctx, key.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to get deleted key if db fails", func(t *testing.T) {
		auth.EXPECT().Check(&types.Operation{Action: types.ActionRead, Resource: types.ResourceKey}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(nil, expectedErr)

		_, err := connector.GetDeleted(ctx, key.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
