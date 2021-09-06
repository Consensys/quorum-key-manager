package eth

import (
	"context"
	"fmt"
	"testing"

	mock3 "github.com/consensys/quorum-key-manager/src/auth/mock"
	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	mock2 "github.com/consensys/quorum-key-manager/src/stores/database/mock"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListEthAccounts(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := fmt.Errorf("error")

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETHAccounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	t.Run("should list ethAccounts successfully", func(t *testing.T) {
		accOne := testutils2.FakeETHAccount()
		accTwo := testutils2.FakeETHAccount()

		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceEthAccount}).Return(nil)
		db.EXPECT().GetAll(gomock.Any()).Return([]*entities.ETHAccount{accOne, accTwo}, nil)

		accAddrs, err := connector.List(ctx)

		assert.NoError(t, err)
		assert.Equal(t, accAddrs, []common.Address{accOne.Address, accTwo.Address})
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceEthAccount}).Return(expectedErr)

		_, err := connector.List(ctx)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to list ethAccounts if db fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceEthAccount}).Return(nil)
		db.EXPECT().GetAll(gomock.Any()).Return(nil, expectedErr)

		_, err := connector.List(ctx)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}

func TestListDeletedEthAccounts(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := fmt.Errorf("error")

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETHAccounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	t.Run("should list deleted ethAccounts successfully", func(t *testing.T) {
		accOne := testutils2.FakeETHAccount()
		accTwo := testutils2.FakeETHAccount()

		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceEthAccount}).Return(nil)
		db.EXPECT().GetAllDeleted(gomock.Any()).Return([]*entities.ETHAccount{accOne, accTwo}, nil)

		accAddrs, err := connector.ListDeleted(ctx)

		assert.NoError(t, err)
		assert.Equal(t, accAddrs, []common.Address{accOne.Address, accTwo.Address})
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceEthAccount}).Return(expectedErr)

		_, err := connector.ListDeleted(ctx)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to list deleted ethAccounts if db fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceEthAccount}).Return(nil)
		db.EXPECT().GetAllDeleted(gomock.Any()).Return(nil, expectedErr)

		_, err := connector.ListDeleted(ctx)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
