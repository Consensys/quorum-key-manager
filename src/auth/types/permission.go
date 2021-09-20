package types

import (
	"fmt"
	"strings"
)

type Permission string

const ReadSecret Permission = "read:secrets"
const WriteSecret Permission = "write:secrets"
const DeleteSecret Permission = "delete:secrets"
const DestroySecret Permission = "destroy:secrets"

const ReadKey Permission = "read:keys"
const WriteKey Permission = "write:keys"
const DeleteKey Permission = "delete:keys"
const DestroyKey Permission = "destroy:keys"
const SignKey Permission = "sign:keys"
const EncryptKey Permission = "encrypt:keys"

const ReadEth Permission = "read:ethereum"
const WriteEth Permission = "write:ethereum"
const DeleteEth Permission = "delete:ethereum"
const DestroyEth Permission = "destroy:ethereum"
const SignEth Permission = "sign:ethereum"
const EncryptEth Permission = "encrypt:ethereum"

const ProxyNode Permission = "proxy:nodes"

func ListPermissions() []Permission {
	return []Permission{
		ReadSecret,
		WriteSecret,
		DeleteSecret,
		DestroySecret,
		ReadKey,
		WriteKey,
		DeleteKey,
		DestroyKey,
		SignKey,
		EncryptKey,
		ReadEth,
		WriteEth,
		DeleteEth,
		DestroyEth,
		SignEth,
		EncryptEth,
		ProxyNode,
	}
}

func ListWildcardPermission(p string) []Permission {
	all := ListPermissions()
	parts := strings.Split(p, ":")
	action, resource := parts[0], parts[1]
	if action == "*" && resource == "*" {
		return all
	}

	var included []Permission
	for _, ip := range all {
		if action == "*" && strings.Contains(string(ip), fmt.Sprintf(":%s", resource)) {
			included = append(included, ip)
		}
		if resource == "*" && strings.Contains(string(ip), fmt.Sprintf("%s:", action)) {
			included = append(included, ip)
		}
	}

	return included
}
