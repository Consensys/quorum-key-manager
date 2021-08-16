package eth1

import (
	"context"
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

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETH1Accounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)

	connector := NewConnector(store, db, nil, logger)

	db.EXPECT().RunInTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, persist func(dbtx database.ETH1Accounts) error) error {
			return persist(db)
		}).AnyTimes()

	t.Run("should update eth1Account successfully", func(t *testing.T) {
		acc := testutils2.FakeETH1Account()
		key := testutils2.FakeKey()
		attributes := testutils2.FakeAttributes()
		acc.Tags = attributes.Tags

		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(acc, nil)

		db.EXPECT().Update(gomock.Any(), acc).Return(acc, nil)

		store.EXPECT().Update(gomock.Any(), acc.Address.Hex(), attributes).Return(key, nil)

		rAcc, err := connector.Update(ctx, acc.Address, attributes)

		assert.NoError(t, err)
		assert.Equal(t, rAcc, acc)
	})

	t.Run("should update key successfully, ignoring not supported error", func(t *testing.T) {
		acc := testutils2.FakeETH1Account()
		attributes := testutils2.FakeAttributes()
		rErr := errors.NotSupportedError("not supported")

		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(acc, nil)

		db.EXPECT().Update(gomock.Any(), acc).Return(acc, nil)

		store.EXPECT().Update(gomock.Any(), acc.Address.Hex(), attributes).Return(nil, rErr)

		_, err := connector.Update(ctx, acc.Address, attributes)

		assert.NoError(t, err)
	})

	t.Run("should fail to update key if key is not found", func(t *testing.T) {
		acc := testutils2.FakeETH1Account()
		attributes := testutils2.FakeAttributes()
		expectedErr := errors.NotFoundError("not found")

		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(nil, expectedErr)

		_, err := connector.Update(ctx, acc.Address, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to update key if db fail to update", func(t *testing.T) {
		acc := testutils2.FakeETH1Account()
		attributes := testutils2.FakeAttributes()
		expectedErr := errors.PostgresError("cannot connect")

		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(acc, nil)

		db.EXPECT().Update(gomock.Any(), acc).Return(nil, expectedErr)

		_, err := connector.Update(ctx, acc.Address, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to update key if store fail to update", func(t *testing.T) {
		acc := testutils2.FakeETH1Account()
		attributes := testutils2.FakeAttributes()
		expectedErr := errors.UnauthorizedError("not authorized")
		acc.Tags = attributes.Tags

		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(acc, nil)

		db.EXPECT().Update(gomock.Any(), acc).Return(acc, nil)

		store.EXPECT().Update(gomock.Any(), acc.Address.Hex(), attributes).Return(nil, expectedErr)

		_, err := connector.Update(ctx, acc.Address, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
