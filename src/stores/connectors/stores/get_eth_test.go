package stores

import (
	"context"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	mock3 "github.com/consensys/quorum-key-manager/src/auth/mock"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	mock2 "github.com/consensys/quorum-key-manager/src/stores/database/mock"
	mock4 "github.com/consensys/quorum-key-manager/src/vaults/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetEthereum(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mock2.NewMockDatabase(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockRoles(ctrl)
	vaults := mock4.NewMockVaults(ctrl)

	connector := NewConnector(auth, db, vaults, logger)

	t.Run("should fail with not found ethereum store successfully", func(t *testing.T) {
		storeName := "not-found-store"
		userInfo := entities.NewWildcardUser()

		auth.EXPECT().UserPermissions(gomock.Any(), userInfo)
		_, err := connector.Ethereum(ctx, storeName, userInfo)

		assert.Error(t, err)
		assert.True(t, errors.IsNotFoundError(err))
	})
}
