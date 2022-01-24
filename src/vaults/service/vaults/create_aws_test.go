package vaults

import (
	"context"
	"testing"

	entities2 "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/auth/mock"
	"github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAWS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := testutils.NewMockLogger(ctrl)
	roles := mock.NewMockRoles(ctrl)
	vault := New(roles, logger)

	ctx := context.Background()
	vaultName := "aws-vault"
	allowedTenantID := "allowed_aws_tenant"
	cfg := &entities.AWSConfig{}
	allowedTenants := []string{allowedTenantID}

	t.Run("should create AWS vault client successfully", func(t *testing.T) {
		userInfo := &entities2.UserInfo{
			Tenant: allowedTenantID,
		}
		err := vault.CreateAWS(ctx, vaultName, cfg, allowedTenants, userInfo)
		assert.NoError(t, err)
	})
}
