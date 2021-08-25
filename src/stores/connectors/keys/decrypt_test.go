package keys

import (
	"context"
	"fmt"
	"testing"

	mock3 "github.com/consensys/quorum-key-manager/src/auth/mock"
	"github.com/consensys/quorum-key-manager/src/auth/types"

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

	key := testutils2.FakeKey()
	expectedErr := fmt.Errorf("error")
	data := []byte("0x123")
	result := []byte("0x456")

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockKeys(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	t.Run("should decrypt data successfully", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionEncrypt, Resource: types.ResourceKey}).Return(nil)
		store.EXPECT().Decrypt(gomock.Any(), key.ID, data).Return(result, nil)

		rResult, err := connector.Decrypt(ctx, key.ID, data)

		assert.NoError(t, err)
		assert.Equal(t, rResult, result)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionEncrypt, Resource: types.ResourceKey}).Return(expectedErr)

		_, err := connector.Decrypt(ctx, key.ID, data)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to decrypt data if decrypt fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&types.Operation{Action: types.ActionEncrypt, Resource: types.ResourceKey}).Return(nil)
		store.EXPECT().Decrypt(gomock.Any(), key.ID, data).Return(nil, expectedErr)

		_, err := connector.Decrypt(ctx, key.ID, data)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
