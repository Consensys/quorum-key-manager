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

func TestCreateHashicorp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := testutils.NewMockLogger(ctrl)
	roles := mock.NewMockRoles(ctrl)
	vault := New(roles, logger)

	ctx := context.Background()
	vaultName := "hashicorp-vault"
	cfg := &entities.HashicorpConfig{}
	allowedTenants := []string{"tenant_id_1"}

	t.Run("should create Hashicorp vault client successfully", func(t *testing.T) {
		userInfo := &entities2.UserInfo{
			Tenant: "tenant_id_1",
		}
		err := vault.CreateHashicorp(ctx, vaultName, cfg, allowedTenants, userInfo)
		assert.NoError(t, err)
	})
}
