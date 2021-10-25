package keys

import (
	"context"
	"fmt"
	"testing"

	mock3 "github.com/consensys/quorum-key-manager/src/auth/mock"
	"github.com/consensys/quorum-key-manager/src/auth/entities"

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

	algo := testutils2.FakeAlgorithm()
	data := []byte("0x123")
	result := []byte("0x456")
	key := testutils2.FakeKey()
	expectedErr := fmt.Errorf("error")

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockKeys(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	t.Run("should sign data successfully", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionSign, Resource: entities.ResourceKey}).Return(nil)
		store.EXPECT().Sign(gomock.Any(), key.ID, data, algo).Return(result, nil)

		rResult, err := connector.Sign(ctx, key.ID, data, algo)

		assert.NoError(t, err)
		assert.Equal(t, rResult, result)
	})

	t.Run("should sign data with key algo successfully", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionSign, Resource: entities.ResourceKey}).Return(nil)
		db.EXPECT().Get(ctx, key.ID).Return(key, nil)
		store.EXPECT().Sign(ctx, key.ID, data, key.Algo).Return(result, nil)

		rResult, err := connector.Sign(ctx, key.ID, data, nil)

		assert.NoError(t, err)
		assert.Equal(t, rResult, result)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionSign, Resource: entities.ResourceKey}).Return(expectedErr)

		_, err := connector.Sign(ctx, key.ID, data, algo)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to sign data if sign fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionSign, Resource: entities.ResourceKey}).Return(nil)
		store.EXPECT().Sign(gomock.Any(), key.ID, data, algo).Return(nil, expectedErr)

		_, err := connector.Sign(ctx, key.ID, data, algo)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to sign data if db fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionSign, Resource: entities.ResourceKey}).Return(nil)
		db.EXPECT().Get(gomock.Any(), key.ID).Return(key, expectedErr)

		_, err := connector.Sign(ctx, key.ID, data, nil)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
