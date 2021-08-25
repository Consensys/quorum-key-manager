package eth1

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

func TestListEth1Accounts(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := fmt.Errorf("error")

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETH1Accounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	t.Run("should list eth1Accounts successfully", func(t *testing.T) {
		accOne := testutils2.FakeETH1Account()
		accTwo := testutils2.FakeETH1Account()

		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceEth1Account}).Return(nil)
		db.EXPECT().GetAll(gomock.Any()).Return([]*entities.ETH1Account{accOne, accTwo}, nil)

		accAddrs, err := connector.List(ctx)

		assert.NoError(t, err)
		assert.Equal(t, accAddrs, []common.Address{accOne.Address, accTwo.Address})
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceEth1Account}).Return(expectedErr)

		_, err := connector.List(ctx)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to list eth1Accounts if db fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceEth1Account}).Return(nil)
		db.EXPECT().GetAll(gomock.Any()).Return(nil, expectedErr)

		_, err := connector.List(ctx)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}

func TestListDeletedEth1Accounts(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := fmt.Errorf("error")

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETH1Accounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	t.Run("should list deleted eth1Accounts successfully", func(t *testing.T) {
		accOne := testutils2.FakeETH1Account()
		accTwo := testutils2.FakeETH1Account()

		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceEth1Account}).Return(nil)
		db.EXPECT().GetAllDeleted(gomock.Any()).Return([]*entities.ETH1Account{accOne, accTwo}, nil)

		accAddrs, err := connector.ListDeleted(ctx)

		assert.NoError(t, err)
		assert.Equal(t, accAddrs, []common.Address{accOne.Address, accTwo.Address})
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceEth1Account}).Return(expectedErr)

		_, err := connector.ListDeleted(ctx)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to list deleted eth1Accounts if db fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceEth1Account}).Return(nil)
		db.EXPECT().GetAllDeleted(gomock.Any()).Return(nil, expectedErr)

		_, err := connector.ListDeleted(ctx)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
