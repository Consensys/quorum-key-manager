package eth

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

func TestDecryptKey(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := []byte("0x123")
	result := []byte("0x456")
	expectedErr := fmt.Errorf("error")
	acc := testutils2.FakeETHAccount()
	key := testutils2.FakeKey()
	key.ID = acc.KeyID

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETHAccounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	t.Run("should decrypt data successfully", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionEncrypt, Resource: entities.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(acc, nil)
		store.EXPECT().Decrypt(gomock.Any(), key.ID, data).Return(result, nil)

		rResult, err := connector.Decrypt(ctx, acc.Address, data)

		assert.NoError(t, err)
		assert.Equal(t, rResult, result)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionEncrypt, Resource: entities.ResourceEthAccount}).Return(expectedErr)

		_, err := connector.Decrypt(ctx, acc.Address, data)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should fail to decrypt data if db fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionEncrypt, Resource: entities.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(nil, expectedErr)

		_, err := connector.Decrypt(ctx, acc.Address, data)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should fail to decrypt data if store fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionEncrypt, Resource: entities.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(acc, nil)
		store.EXPECT().Decrypt(gomock.Any(), key.ID, data).Return(nil, expectedErr)

		_, err := connector.Decrypt(ctx, acc.Address, data)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
