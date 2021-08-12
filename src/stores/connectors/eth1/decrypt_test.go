package eth1

import (
	"context"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	mock2 "github.com/consensys/quorum-key-manager/src/stores/database/mock"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDecryptKey(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETH1Accounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)

	data := []byte("0x123")
	result := []byte("0x456")
	connector := NewConnector(store, db, logger)

	t.Run("should decrypt data successfully", func(t *testing.T) {
		acc := testutils2.FakeETH1Account()
		key := testutils2.FakeKey()
		attributes := testutils2.FakeAttributes()
		key.ID = acc.KeyID
		acc.Tags = attributes.Tags

		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(acc, nil)

		store.EXPECT().Decrypt(gomock.Any(), key.ID, data).Return(result, nil)

		rResult, err := connector.Decrypt(ctx, acc.Address, data)

		assert.NoError(t, err)
		assert.Equal(t, rResult, result)
	})

	t.Run("should fail to decrypt data if db fails", func(t *testing.T) {
		expectedErr := errors.PostgresError("cannot connect")
		acc := testutils2.FakeETH1Account()
		key := testutils2.FakeKey()
		attributes := testutils2.FakeAttributes()
		key.ID = acc.KeyID
		acc.Tags = attributes.Tags

		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(nil, expectedErr)

		_, err := connector.Decrypt(ctx, acc.Address, data)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should fail to decrypt data if store fails", func(t *testing.T) {
		expectedErr := errors.UnauthorizedError("unauthorized")
		acc := testutils2.FakeETH1Account()
		key := testutils2.FakeKey()
		attributes := testutils2.FakeAttributes()
		key.ID = acc.KeyID
		acc.Tags = attributes.Tags

		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(acc, nil)

		store.EXPECT().Decrypt(gomock.Any(), key.ID, data).Return(nil, expectedErr)

		_, err := connector.Decrypt(ctx, acc.Address, data)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
