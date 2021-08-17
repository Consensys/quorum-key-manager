package eth1

import (
	"context"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"
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

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETH1Accounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)

	connector := NewConnector(store, db, nil, logger)

	t.Run("should list eth1Accounts successfully", func(t *testing.T) {
		accOne := testutils2.FakeETH1Account()
		accTwo := testutils2.FakeETH1Account()

		db.EXPECT().GetAll(gomock.Any()).Return([]*entities.ETH1Account{accOne, accTwo}, nil)

		accAddrs, err := connector.List(ctx)

		assert.NoError(t, err)
		assert.Equal(t, accAddrs, []common.Address{accOne.Address, accTwo.Address})
	})

	t.Run("should fail to list eth1Accounts if db fails", func(t *testing.T) {
		expectedErr := errors.PostgresError("cannot connect")

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

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETH1Accounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)

	connector := NewConnector(store, db, nil, logger)

	t.Run("should list deleted eth1Accounts successfully", func(t *testing.T) {
		accOne := testutils2.FakeETH1Account()
		accTwo := testutils2.FakeETH1Account()

		db.EXPECT().GetAllDeleted(gomock.Any()).Return([]*entities.ETH1Account{accOne, accTwo}, nil)

		accAddrs, err := connector.ListDeleted(ctx)

		assert.NoError(t, err)
		assert.Equal(t, accAddrs, []common.Address{accOne.Address, accTwo.Address})
	})

	t.Run("should fail to list deleted eth1Accounts if db fails", func(t *testing.T) {
		expectedErr := errors.PostgresError("cannot connect")

		db.EXPECT().GetAllDeleted(gomock.Any()).Return(nil, expectedErr)

		_, err := connector.ListDeleted(ctx)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
