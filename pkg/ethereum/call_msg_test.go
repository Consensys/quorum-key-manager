package ethereum

import (
	"encoding/json"
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertCallMsgEquals(t *testing.T, expectedMsg, msg *CallMsg) {
	assert.Equal(t, expectedMsg.From, msg.From, "From should be correct")
	assert.Equal(t, expectedMsg.To, msg.To, "To should be correct")
	assert.Equal(t, expectedMsg.Gas, msg.Gas, "Gas should be correct")
	assert.Equal(t, expectedMsg.GasPrice, msg.GasPrice, "GasPrice should be correct")
	assert.Equal(t, expectedMsg.Value, msg.Value, "Value should be correct")
	assert.Equal(t, expectedMsg.Data, msg.Data, "Data should be correct")
}

func TestCallMsg(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the response
		body []byte

		expectedCallMsg CallMsg
		expectedErrMsg  string
	}{
		{
			desc: "all fields",
			body: []byte(`{"from":"0xc94770007dda54cf92009bff0de90c06f603a09f","to":"0xfe3b557e8fb62b89f4916b721be55ceb828dbd73","gas":"0x5208","gasPrice":"0x3e8","value":"0x1","data":"0xabcdef"}`),
			expectedCallMsg: CallMsg{
				From:     func(addr ethcommon.Address) *ethcommon.Address { return &addr }(ethcommon.HexToAddress("0xc94770007dda54cf92009bff0de90c06f603a09f")),
				To:       func(addr ethcommon.Address) *ethcommon.Address { return &addr }(ethcommon.HexToAddress("0xfe3b557e8fb62b89f4916b721be55ceb828dbd73")),
				Gas:      func(i uint64) *uint64 { return &i }(21000),
				GasPrice: big.NewInt(1000),
				Value:    big.NewInt(1),
				Data:     func(b []byte) *[]byte { return &b }(ethcommon.FromHex("0xabcdef")),
			},
		},
		{
			desc: "partial fields",
			body: []byte(`{"from":"0xc94770007dda54cf92009bff0de90c06f603a09f","to":"0xfe3b557e8fb62b89f4916b721be55ceb828dbd73","value":"0x1"}`),
			expectedCallMsg: CallMsg{
				From:  func(addr ethcommon.Address) *ethcommon.Address { return &addr }(ethcommon.HexToAddress("0xc94770007dda54cf92009bff0de90c06f603a09f")),
				To:    func(addr ethcommon.Address) *ethcommon.Address { return &addr }(ethcommon.HexToAddress("0xfe3b557e8fb62b89f4916b721be55ceb828dbd73")),
				Value: big.NewInt(1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			msg := new(CallMsg)
			err := json.Unmarshal(tt.body, msg)
			if tt.expectedErrMsg == "" {
				require.NoError(t, err, "Unmarshal must not fail")
				assertCallMsgEquals(t, &tt.expectedCallMsg, msg)
				b, err := json.Marshal(msg)
				require.NoError(t, err, "Marshal must not fail")
				assert.Equal(t, tt.body, b, "Marshal body should match")
			}
		})
	}
}
