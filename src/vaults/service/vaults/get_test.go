package vaults

import (
	"context"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	entities2 "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/auth/mock"
	"github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetVault(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := testutils.NewMockLogger(ctrl)
	roles := mock.NewMockRoles(ctrl)
	vault := New(roles, logger)

	ctx := context.Background()
	vaultName := "vault-id"
	allowedTenantID := "allowed_tenant"
	cfg := &entities.AWSConfig{}
	allowedTenants := []string{allowedTenantID}

	t.Run("should create and get vault client successfully", func(t *testing.T) {
		userInfo := &entities2.UserInfo{
			Tenant: allowedTenantID,
		}
		err := vault.CreateAWS(ctx, vaultName, cfg, allowedTenants, userInfo)
		assert.NoError(t, err)

		roles.EXPECT().UserPermissions(ctx, userInfo).Return([]entities2.Permission{})
		vault, err := vault.Get(ctx, vaultName, userInfo)
		assert.NoError(t, err)
		assert.NotNil(t, vault)
	})

	t.Run("should fail to get vault client when does not exists", func(t *testing.T) {
		userInfo := &entities2.UserInfo{
			Tenant: allowedTenantID,
		}

		roles.EXPECT().UserPermissions(ctx, userInfo).Return([]entities2.Permission{})
		_, err := vault.Get(ctx, "not-existing-vault", userInfo)
		assert.True(t, errors.IsNotFoundError(err))
	})

	t.Run("should fail to get vault client if tenant is not allowed", func(t *testing.T) {
		vaultName2 := "aws-vault-2"
		userInfo := &entities2.UserInfo{
			Tenant: "invalid_tenant_id",
		}
		err := vault.CreateAWS(ctx, vaultName2, cfg, allowedTenants, userInfo)
		assert.NoError(t, err)

		roles.EXPECT().UserPermissions(ctx, userInfo).Return([]entities2.Permission{})
		_, err = vault.Get(ctx, vaultName2, userInfo)
		assert.True(t, errors.IsNotFoundError(err))
	})
}
