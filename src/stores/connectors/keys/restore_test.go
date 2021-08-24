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

func TestRestoreKey(t *testing.T) {
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

	db.EXPECT().RunInTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, persist func(dbtx database.Keys) error) error {
			return persist(db)
		}).AnyTimes()

	t.Run("should restore key successfully", func(t *testing.T) {
		auth.EXPECT().Check(&types.Operation{Action: types.ActionDelete, Resource: types.ResourceKey}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Restore(gomock.Any(), key.ID).Return(nil)
		store.EXPECT().Restore(gomock.Any(), key.ID).Return(nil)

		err := connector.Restore(ctx, key.ID)

		assert.NoError(t, err)
	})

	t.Run("should restore key successfully, ignoring not supported error", func(t *testing.T) {
		rErr := errors.NotSupportedError("not supported")

		auth.EXPECT().Check(&types.Operation{Action: types.ActionDelete, Resource: types.ResourceKey}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Restore(gomock.Any(), key.ID).Return(nil)
		store.EXPECT().Restore(gomock.Any(), key.ID).Return(rErr)

		err := connector.Restore(ctx, key.ID)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().Check(&types.Operation{Action: types.ActionDelete, Resource: types.ResourceKey}).Return(expectedErr)

		err := connector.Restore(ctx, key.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to restore key if key is not deleted", func(t *testing.T) {
		auth.EXPECT().Check(&types.Operation{Action: types.ActionDelete, Resource: types.ResourceKey}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(key, expectedErr)

		err := connector.Restore(ctx, key.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to restore key if db fail to restore", func(t *testing.T) {
		auth.EXPECT().Check(&types.Operation{Action: types.ActionDelete, Resource: types.ResourceKey}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Restore(gomock.Any(), key.ID).Return(expectedErr)

		err := connector.Restore(ctx, key.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to restore key if store fail to restore", func(t *testing.T) {
		auth.EXPECT().Check(&types.Operation{Action: types.ActionDelete, Resource: types.ResourceKey}).Return(nil)
		db.EXPECT().GetDeleted(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Restore(gomock.Any(), key.ID).Return(nil)
		store.EXPECT().Restore(gomock.Any(), key.ID).Return(expectedErr)

		err := connector.Restore(ctx, key.ID)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
