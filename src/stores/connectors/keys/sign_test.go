package keys

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

func TestSignKey(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockKeys(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	algo := testutils2.FakeAlgorithm()
	data := []byte("0x123")
	result := []byte("0x456")

	connector := NewConnector(store, db, logger)

	t.Run("should sign data successfully", func(t *testing.T) {
		key := testutils2.FakeKey()

		store.EXPECT().Sign(gomock.Any(), key.ID, data, algo).Return(result, nil)

		rResult, err := connector.Sign(ctx, key.ID, data, algo)

		assert.NoError(t, err)
		assert.Equal(t, rResult, result)
	})

	t.Run("should sign data with key algo successfully", func(t *testing.T) {
		key := testutils2.FakeKey()

		db.EXPECT().Get(gomock.Any(), key.ID).Return(key, nil)
		store.EXPECT().Sign(gomock.Any(), key.ID, data, key.Algo).Return(result, nil)

		rResult, err := connector.Sign(ctx, key.ID, data, nil)

		assert.NoError(t, err)
		assert.Equal(t, rResult, result)
	})

	t.Run("should fail to sign data if sign fails", func(t *testing.T) {
		key := testutils2.FakeKey()
		expectedErr := errors.UnauthorizedError("not authorized")

		store.EXPECT().Sign(gomock.Any(), key.ID, data, algo).Return(nil, expectedErr)

		_, err := connector.Sign(ctx, key.ID, data, algo)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to sign data if db fails", func(t *testing.T) {
		key := testutils2.FakeKey()
		expectedErr := errors.PostgresError("cannot connect")

		db.EXPECT().Get(gomock.Any(), key.ID).Return(key, expectedErr)

		_, err := connector.Sign(ctx, key.ID, data, nil)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
