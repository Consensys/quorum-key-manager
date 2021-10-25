package eth

import (
	"context"
	"fmt"
	"testing"

	mock3 "github.com/consensys/quorum-key-manager/src/auth/mock"
	"github.com/consensys/quorum-key-manager/src/auth/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	mock2 "github.com/consensys/quorum-key-manager/src/stores/database/mock"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRestoreEthAccount(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := fmt.Errorf("error")
	acc := testutils2.FakeETHAccount()

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETHAccounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	db.EXPECT().RunInTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, persist func(dbtx database.ETHAccounts) error) error {
			return persist(db)
		}).AnyTimes()

	t.Run("should restore ethAccount successfully", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceEthAccount}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(nil, errors.NotFoundError("error"))
		db.EXPECT().GetDeleted(gomock.Any(), acc.Address.Hex()).Return(acc, nil)
		db.EXPECT().Restore(gomock.Any(), acc.Address.Hex()).Return(nil)
		store.EXPECT().Restore(gomock.Any(), acc.KeyID).Return(nil)

		err := connector.Restore(ctx, acc.Address)

		assert.NoError(t, err)
	})

	t.Run("should restore ethAccount successfully, ignoring not supported error", func(t *testing.T) {
		rErr := errors.NotSupportedError("not supported")
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceEthAccount}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(nil, errors.NotFoundError(""))
		db.EXPECT().GetDeleted(gomock.Any(), acc.Address.Hex()).Return(acc, nil)
		db.EXPECT().Restore(gomock.Any(), acc.Address.Hex()).Return(nil)
		store.EXPECT().Restore(gomock.Any(), acc.KeyID).Return(rErr)

		err := connector.Restore(ctx, acc.Address)

		assert.NoError(t, err)
	})

	t.Run("should be idempotent if ethAccount already exists", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceEthAccount}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(nil, nil)

		err := connector.Restore(ctx, acc.Address)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceEthAccount}).Return(expectedErr)

		err := connector.Restore(ctx, acc.Address)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to restore ethAccount if ethAccount is not yet deleted", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceEthAccount}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(nil, errors.NotFoundError(""))
		db.EXPECT().GetDeleted(gomock.Any(), acc.Address.Hex()).Return(nil, expectedErr)

		err := connector.Restore(ctx, acc.Address)

		assert.Error(t, err)
	})

	t.Run("should fail to restore ethAccount if db fails to restore", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceEthAccount}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(nil, errors.NotFoundError(""))
		db.EXPECT().GetDeleted(gomock.Any(), acc.Address.Hex()).Return(acc, nil)
		db.EXPECT().Restore(gomock.Any(), acc.Address.Hex()).Return(expectedErr)

		err := connector.Restore(ctx, acc.Address)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to restore ethAccount if store fails to restore", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceEthAccount}).Return(nil)
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(nil, errors.NotFoundError(""))
		db.EXPECT().GetDeleted(gomock.Any(), acc.Address.Hex()).Return(acc, nil)
		db.EXPECT().Restore(gomock.Any(), acc.Address.Hex()).Return(nil)
		store.EXPECT().Restore(gomock.Any(), acc.KeyID).Return(expectedErr)

		err := connector.Restore(ctx, acc.Address)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
