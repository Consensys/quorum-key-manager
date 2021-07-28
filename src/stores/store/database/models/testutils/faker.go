package testutils

import (
	"time"

	"github.com/consensys/quorum-key-manager/src/stores/store/database/models"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities/testutils"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func FakeETH1Account() *models.ETH1Account {
	return &models.ETH1Account{
		KeyID:               "my-account",
		Address:             "0x83a0254be47813BBff771F4562744676C4e793F0",
		PublicKey:           hexutil.MustDecode("0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2"),
		CompressedPublicKey: hexutil.MustDecode("0x02555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f"),
		Tags:                testutils.FakeTags(),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}
