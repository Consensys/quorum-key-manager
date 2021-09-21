package header

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOverride(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the response
		dst, overrides, expectedDst map[string][]string
	}{
		{
			desc: "no overlapping header",
			dst: map[string][]string{
				"HEADER1": {"dst-value1.1"},
				"HEADER2": {"dst-value1.1", "dst-value1.2"},
			},
			overrides: map[string][]string{
				"HEADER3": {"src-value3.1"},
			},
			expectedDst: map[string][]string{
				"HEADER1": {"dst-value1.1"},
				"HEADER2": {"dst-value1.1", "dst-value1.2"},
				"HEADER3": {"src-value3.1"},
			},
		},
		{
			desc: "overlapping header",
			dst: map[string][]string{
				"HEADER1": {"dst-value1.1"},
				"HEADER2": {"dst-value1.1", "dst-value1.2"},
			},
			overrides: map[string][]string{
				"HEADER2": {"src-value2.1"},
			},
			expectedDst: map[string][]string{
				"HEADER1": {"dst-value1.1"},
				"HEADER2": {"src-value2.1"},
			},
		},
		{
			desc: "deleting header",
			dst: map[string][]string{
				"HEADER1": {"dst-value1.1"},
				"HEADER2": {"dst-value1.1", "dst-value1.2"},
			},
			overrides: map[string][]string{
				"HEADER2": {},
			},
			expectedDst: map[string][]string{
				"HEADER1": {"dst-value1.1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dst := FromMap(tt.dst)
			expectedDst := FromMap(tt.expectedDst)
			_ = Override(tt.overrides)(dst)
			assert.Equal(t, expectedDst, dst, "Override should be correct")
		})
	}
}
