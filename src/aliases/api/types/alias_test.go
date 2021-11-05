package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshallingJSONOK(t *testing.T) {
	cases := []string{
		`{
		"type": "string",
		"value": "key1"
	}`,

		`{
		"type": "array",
		"value": ["key2"]
	}`,
	}

	for _, c := range cases {
		var av AliasRequest
		err := json.Unmarshal([]byte(c), &av)
		assert.NoError(t, err)
		b, err := json.Marshal(av)
		assert.NoError(t, err)
		assert.JSONEq(t, c, string(b))
	}
}

func TestUnmarshalJSONError(t *testing.T) {
	cases := []struct {
		json string
		req  AliasRequest
	}{
		{
			json: `{
		"type": "arrayzzz",
		"value": ["key3"]
	}`,
			req: AliasRequest{AliasValue: AliasValue{RawKind: "arrayzzz", RawValue: []string{"key3"}}},
		},
	}
	for _, c := range cases {
		var av AliasRequest
		err := json.Unmarshal([]byte(c.json), &av)
		assert.Error(t, err)
		_, err = json.Marshal(c.req)
		assert.Error(t, err)
	}
}
