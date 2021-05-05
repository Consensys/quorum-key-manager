package ethereum

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockNumber(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the response
		body []byte

		expectedBlockNumber BlockNumber
	}{
		{
			desc:                "pending",
			body:                []byte(`"pending"`),
			expectedBlockNumber: -2,
		},
		{
			desc:                "latest",
			body:                []byte(`"latest"`),
			expectedBlockNumber: -1,
		},
		{
			desc:                "earliest",
			body:                []byte(`"earliest"`),
			expectedBlockNumber: 0,
		},
		{
			desc:                "number",
			body:                []byte(`"0xf"`),
			expectedBlockNumber: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			bn := new(BlockNumber)
			err := json.Unmarshal(tt.body, bn)
			require.NoError(t, err, "Unmarshal must not fail")
			assert.Equal(t, tt.expectedBlockNumber, *bn, "Unmarshal should be valid")
			b, err := json.Marshal(bn)
			require.NoError(t, err, "Marshal must not fail")
			assert.Equal(t, tt.body, b, "Marshal body should match")
		})
	}
}
