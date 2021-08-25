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

func TestCreateKey(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := testutils2.FakeKey()
	expectedErr := fmt.Errorf("error")
	attributes := testutils2.FakeAttributes()

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockKeys(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	t.Run("should create key successfully", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceKey}).Return(nil)
		store.EXPECT().Create(gomock.Any(), key.ID, key.Algo, attributes).Return(key, nil)
		db.EXPECT().Add(gomock.Any(), key).Return(key, nil)

		rKey, err := connector.Create(ctx, key.ID, key.Algo, attributes)

		assert.NoError(t, err)
		assert.Equal(t, rKey, key)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceKey}).Return(expectedErr)

		_, err := connector.Create(ctx, key.ID, key.Algo, attributes)

		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to delete key if store fail to create", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceKey}).Return(nil)
		store.EXPECT().Create(gomock.Any(), key.ID, key.Algo, attributes).Return(nil, expectedErr)

		_, err := connector.Create(ctx, key.ID, key.Algo, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to create key if db fail to add", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceKey}).Return(nil)
		store.EXPECT().Create(gomock.Any(), key.ID, key.Algo, attributes).Return(key, nil)
		db.EXPECT().Add(gomock.Any(), key).Return(nil, expectedErr)

		_, err := connector.Create(ctx, key.ID, key.Algo, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
