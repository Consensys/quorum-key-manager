package eth1

import (
	"context"
	"testing"

	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	mock2 "github.com/consensys/quorum-key-manager/src/stores/database/mock"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecover(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETH1Accounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)

	connector := NewConnector(store, db, logger)
	acc := testutils2.FakeETH1Account()
	acc.Address = common.HexToAddress("0x6436Bd740B732b90a9f7bc3065d6c3eDa57D9785")

	t.Run("should recover account successfully", func(t *testing.T) {
		data := crypto.Keccak256([]byte("my data to sign"))
		ecdsaSignature := hexutil.MustDecode("0x314EDF887EECB3C4BA7C90F9BD03D1044BC53EB2CADCE8C1E056768ACF8904372B8759BBCA88341BF074BB0595E6A19B7167BE6DA6D5687E81892E10B349D6FE01")

		address, err := connector.ECRecover(ctx, data, ecdsaSignature)

		require.NoError(t, err)
		assert.Equal(t, address, acc.Address)
	})

	t.Run("should fail to recover account if signature is invalid", func(t *testing.T) {
		data := crypto.Keccak256([]byte("my data to sign"))
		ecdsaSignature := hexutil.MustDecode("0x4EDF887EECB3C4BA7C90F9BD03D1044BC53EB2CADCE8C1E056768ACF8904372B8759BBCA88341BF074BB0595E6A19B7167BE6DA6D5687E81892E10B349D6FE01")

		_, err := connector.ECRecover(ctx, data, ecdsaSignature)

		assert.Error(t, err)
	})
}
