package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListWildcardPermission(t *testing.T) {
	list := ListWildcardPermission("*:*")
	assert.Equal(t, list, ListPermissions())

	list = ListWildcardPermission("read:*")
	assert.Equal(t, list, []Permission{ReadSecret, ReadKey, ReadEth1})

	list = ListWildcardPermission("*:eth1accounts")
	assert.Equal(t, list, []Permission{ReadEth1, WriteEth1, DeleteEth1, DestroyEth1, SignEth1, EncryptEth1})
}
