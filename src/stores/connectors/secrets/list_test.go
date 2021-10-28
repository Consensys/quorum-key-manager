package secrets

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
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListSecret(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := fmt.Errorf("error")

	store := mock.NewMockSecretStore(ctrl)
	db := mock2.NewMockSecrets(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)
	connector := NewConnector(store, db, auth, logger)

	t.Run("should list secrets successfully", func(t *testing.T) {
		secretOne := testutils2.FakeSecret()
		secretTwo := testutils2.FakeSecret()
		limit := uint64(2)
		offset := uint64(4)

		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceSecret}).Return(nil)
		db.EXPECT().SearchIDs(gomock.Any(), false, limit, offset).Return([]string{secretOne.ID, secretTwo.ID}, nil)

		secretIDs, err := connector.List(ctx, limit, offset)

		assert.NoError(t, err)
		assert.Equal(t, secretIDs, []string{secretOne.ID, secretTwo.ID})
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceSecret}).Return(expectedErr)

		_, err := connector.List(ctx, uint64(0), uint64(0))

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to list deleted secret if db fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceSecret}).Return(nil)
		db.EXPECT().SearchIDs(gomock.Any(), false, uint64(0), uint64(0)).Return(nil, expectedErr)

		_, err := connector.List(ctx, uint64(0), uint64(0))

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}

func TestListDeletedSecret(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := fmt.Errorf("error")

	store := mock.NewMockSecretStore(ctrl)
	db := mock2.NewMockSecrets(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	t.Run("should list deleted secret successfully", func(t *testing.T) {
		secretOne := testutils2.FakeSecret()
		secretTwo := testutils2.FakeSecret()
		limit := uint64(2)
		offset := uint64(4)

		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceSecret}).Return(nil)
		db.EXPECT().SearchIDs(gomock.Any(), true, limit, offset).Return([]string{secretOne.ID, secretTwo.ID}, nil)

		secretIDs, err := connector.ListDeleted(ctx, limit, offset)

		assert.NoError(t, err)
		assert.Equal(t, secretIDs, []string{secretOne.ID, secretTwo.ID})
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceSecret}).Return(expectedErr)

		_, err := connector.ListDeleted(ctx, uint64(0), uint64(0))

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to list deleted secret if db fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceSecret}).Return(nil)
		db.EXPECT().SearchIDs(gomock.Any(), true, uint64(0), uint64(0)).Return(nil, expectedErr)

		_, err := connector.ListDeleted(ctx, uint64(0), uint64(0))

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
