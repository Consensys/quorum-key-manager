package eth

import (
	"context"
	"fmt"
	"testing"

	"github.com/consensys/quorum-key-manager/src/stores/database/models"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	mock3 "github.com/consensys/quorum-key-manager/src/auth/mock"
	"github.com/consensys/quorum-key-manager/src/auth/entities"

	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	mock2 "github.com/consensys/quorum-key-manager/src/stores/database/mock"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := fmt.Errorf("error")
	acc := testutils2.FakeETHAccount()
	key := testutils2.FakeKey()
	attributes := testutils2.FakeAttributes()
	key.ID = acc.KeyID
	acc.Tags = attributes.Tags
	privKey := []byte("0xABCD")

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETHAccounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	t.Run("should import eth account successfully", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionWrite, Resource: entities.ResourceEthAccount}).Return(nil)
		store.EXPECT().Import(gomock.Any(), key.ID, privKey, ethAlgo, attributes).Return(key, nil)
		db.EXPECT().Add(gomock.Any(), models.NewETHAccountFromKey(key, attributes)).Return(acc, nil)

		rAcc, err := connector.Import(ctx, key.ID, privKey, attributes)

		assert.NoError(t, err)
		assert.Equal(t, rAcc, acc)
	})

	t.Run("should import eth account successfully if it already exists in the vault", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionWrite, Resource: entities.ResourceEthAccount}).Return(nil)
		store.EXPECT().Import(gomock.Any(), key.ID, privKey, ethAlgo, attributes).Return(nil, errors.AlreadyExistsError("error"))
		store.EXPECT().Get(gomock.Any(), key.ID).Return(key, nil)
		db.EXPECT().Add(gomock.Any(), models.NewETHAccountFromKey(key, attributes)).Return(acc, nil)

		rAcc, err := connector.Import(ctx, key.ID, privKey, attributes)

		assert.NoError(t, err)
		assert.Equal(t, rAcc, acc)
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionWrite, Resource: entities.ResourceEthAccount}).Return(expectedErr)

		_, err := connector.Import(ctx, key.ID, privKey, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to create ethAccount if store fail to create", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionWrite, Resource: entities.ResourceEthAccount}).Return(nil)
		store.EXPECT().Import(gomock.Any(), key.ID, privKey, ethAlgo, attributes).Return(nil, expectedErr)

		_, err := connector.Import(ctx, key.ID, privKey, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to create ethAccount if db fail to add", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionWrite, Resource: entities.ResourceEthAccount}).Return(nil)
		store.EXPECT().Import(gomock.Any(), key.ID, privKey, ethAlgo, attributes).Return(key, nil)
		db.EXPECT().Add(gomock.Any(), models.NewETHAccountFromKey(key, attributes)).Return(acc, expectedErr)

		_, err := connector.Import(ctx, key.ID, privKey, attributes)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
