package utils

import (
	"testing"

	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestKeysVerifyMessage_ecdsa256k1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := testutils.NewMockLogger(ctrl)

	connector := New(logger)
	key.

	t.Run("should verify message successfully with recID 27", func(t *testing.T) {
		data := crypto.Keccak256([]byte("my data to sign"))
		ecdsaSignature := hexutil.MustDecode("0x314EDF887EECB3C4BA7C90F9BD03D1044BC53EB2CADCE8C1E056768ACF8904372B8759BBCA88341BF074BB0595E6A19B7167BE6DA6D5687E81892E10B349D6FE1B")
		err := connector.Verify(key.PublicKey, data, ecdsaSignature, key.Algo)

		assert.NoError(t, err)
	})
}
