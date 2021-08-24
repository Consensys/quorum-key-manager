package keys

import (
	"context"
	"fmt"
	mock3 "github.com/consensys/quorum-key-manager/src/auth/mock"
	"github.com/consensys/quorum-key-manager/src/auth/types"
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

func TestUpdateKey(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := testutils2.FakeKey()
	attributes := testutils2.FakeAttributes()
	expectedErr := fmt.Errorf("error")

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockKeys(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	db.EXPECT().RunInTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, persist func(dbtx database.Keys) error) error {
			return persist(db)
		}).AnyTimes()

	t.Run("should update key successfully", func(t *testing.T) {
		updatedKey := testutils2.FakeKey()
		updatedKey.Tags = attributes.Tags

		auth.EXPECT().Check(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceKey}).Return(nil)
		db.EXPECT().Get(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Update(gomock.Any(), key).Return(updatedKey, nil)
		store.EXPECT().Update(gomock.Any(), key.ID, attributes).Return(updatedKey, nil)

		rKey, err := connector.Update(ctx, key.ID, attributes)

		assert.NoError(t, err)
		assert.Equal(t, rKey, updatedKey)
	})

	t.Run("should update key successfully, ignoring not supported error", func(t *testing.T) {
		rErr := errors.NotSupportedError("not supported")

		auth.EXPECT().Check(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceKey}).Return(nil)
		db.EXPECT().Get(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Update(gomock.Any(), key).Return(key, nil)
		store.EXPECT().Update(gomock.Any(), key.ID, attributes).Return(nil, rErr)

		_, err := connector.Update(ctx, key.ID, attributes)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().Check(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceKey}).Return(expectedErr)

		_, err := connector.Update(ctx, key.ID, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to update key if key is not found", func(t *testing.T) {
		auth.EXPECT().Check(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceKey}).Return(nil)
		db.EXPECT().Get(gomock.Any(), key.ID).Return(key, expectedErr)

		_, err := connector.Update(ctx, key.ID, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to update key if db fail to update", func(t *testing.T) {
		auth.EXPECT().Check(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceKey}).Return(nil)
		db.EXPECT().Get(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Update(gomock.Any(), key).Return(nil, expectedErr)

		_, err := connector.Update(ctx, key.ID, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to update key if store fail to update", func(t *testing.T) {
		auth.EXPECT().Check(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceKey}).Return(nil)
		db.EXPECT().Get(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Update(gomock.Any(), key).Return(key, nil)
		store.EXPECT().Update(gomock.Any(), key.ID, attributes).Return(nil, expectedErr)

		_, err := connector.Update(ctx, key.ID, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
