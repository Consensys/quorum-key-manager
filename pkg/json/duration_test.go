package json

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the request
		bytes            []byte
		expectedDuration Duration
	}{
		{
			desc:             "int",
			bytes:            []byte(`20`),
			expectedDuration: Duration{time.Duration(20)},
		},
		{
			desc:             "string",
			bytes:            []byte(`"15s30ns"`),
			expectedDuration: Duration{time.Duration(15000000030)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dur := new(Duration)
			err := json.Unmarshal(tt.bytes, dur)
			require.NoError(t, err, "Unmarshal should not error")
			assert.Equal(t, tt.expectedDuration, *dur, "Duration should be correct")
		})
	}
}
