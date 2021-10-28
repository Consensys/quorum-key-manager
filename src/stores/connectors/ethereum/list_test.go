package eth

import (
	"context"
	"fmt"
	"testing"

	"github.com/consensys/quorum-key-manager/src/auth/entities"
	mock3 "github.com/consensys/quorum-key-manager/src/auth/mock"

	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	mock2 "github.com/consensys/quorum-key-manager/src/stores/database/mock"
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
		limit := uint64(2)
		offset := uint64(4)

		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceEthAccount}).Return(nil)
		db.EXPECT().SearchAddresses(gomock.Any(), false, limit, offset).Return([]string{accOne.Address.String(), accTwo.Address.String()}, nil)

		accAddrs, err := connector.List(ctx, limit, offset)

		assert.NoError(t, err)
		assert.Equal(t, accAddrs, []common.Address{accOne.Address, accTwo.Address})
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceEthAccount}).Return(expectedErr)

		_, err := connector.List(ctx, 0, 0)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to list ethAccounts if db fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceEthAccount}).Return(nil)
		db.EXPECT().SearchAddresses(gomock.Any(), false, uint64(0), uint64(0)).Return(nil, expectedErr)

		_, err := connector.List(ctx, 0, 0)

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
		limit := uint64(2)
		offset := uint64(4)

		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceEthAccount}).Return(nil)
		db.EXPECT().SearchAddresses(gomock.Any(), true, limit, offset).Return([]string{accOne.Address.String(), accTwo.Address.String()}, nil)

		accAddrs, err := connector.ListDeleted(ctx, limit, offset)

		assert.NoError(t, err)
		assert.Equal(t, accAddrs, []common.Address{accOne.Address, accTwo.Address})
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceEthAccount}).Return(expectedErr)

		_, err := connector.ListDeleted(ctx, uint64(0), uint64(0))

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to list deleted ethAccounts if db fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceEthAccount}).Return(nil)
		db.EXPECT().SearchAddresses(gomock.Any(), true, uint64(0), uint64(0)).Return(nil, expectedErr)

		_, err := connector.ListDeleted(ctx, uint64(0), uint64(0))

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
