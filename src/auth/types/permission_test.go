package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListWildcardPermission(t *testing.T) {
	list := ListWildcardPermission("*:*")
	assert.Equal(t, list, ListPermissions())

	list = ListWildcardPermission("read:*")
	assert.Equal(t, list, []Permission{ReadSecret, ReadKey, ReadEth})

	list = ListWildcardPermission("*:ethaccounts")
	assert.Equal(t, list, []Permission{ReadEth, WriteEth, DeleteEth, DestroyEth, SignEth, EncryptEth})
}
