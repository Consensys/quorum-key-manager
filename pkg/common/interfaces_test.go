package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterfaceToObject(t *testing.T) {
	t.Run("parse interface into array of string successfully", func(t *testing.T) {
		var res []string
		err := InterfaceToObject([]interface{}{"a", "b", "c"}, &res)
		assert.NoError(t, err)
		assert.Equal(t, res[0], "a")
	})
	
	t.Run("fail to parse no corresponding interface types", func(t *testing.T) {
		var res []int
		err := InterfaceToObject([]interface{}{"a", "b", "c"}, &res)
		assert.Error(t, err)
	})
}
