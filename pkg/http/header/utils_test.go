package header

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOveride(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the response
		dst, overides, expectedDst map[string][]string
	}{
		{
			desc: "no overlapping header",
			dst: map[string][]string{
				"HEADER1": []string{"dst-value1.1"},
				"HEADER2": []string{"dst-value1.1", "dst-value1.2"},
			},
			overides: map[string][]string{
				"HEADER3": []string{"src-value3.1"},
			},
			expectedDst: map[string][]string{
				"HEADER1": []string{"dst-value1.1"},
				"HEADER2": []string{"dst-value1.1", "dst-value1.2"},
				"HEADER3": []string{"src-value3.1"},
			},
		},
		{
			desc: "overlapping header",
			dst: map[string][]string{
				"HEADER1": []string{"dst-value1.1"},
				"HEADER2": []string{"dst-value1.1", "dst-value1.2"},
			},
			overides: map[string][]string{
				"HEADER2": []string{"src-value2.1"},
			},
			expectedDst: map[string][]string{
				"HEADER1": []string{"dst-value1.1"},
				"HEADER2": []string{"dst-value1.1", "dst-value1.2", "src-value2.1"},
			},
		},
		{
			desc: "deleting header",
			dst: map[string][]string{
				"HEADER1": []string{"dst-value1.1"},
				"HEADER2": []string{"dst-value1.1", "dst-value1.2"},
			},
			overides: map[string][]string{
				"HEADER2": []string{},
			},
			expectedDst: map[string][]string{
				"HEADER1": []string{"dst-value1.1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dst := FromMap(tt.dst)
			expectedDst := FromMap(tt.expectedDst)
			Overide(dst, tt.overides)
			assert.Equal(t, expectedDst, dst, "Overide should be correct")
		})
	}
}
