package types

import (
	"fmt"
	"strings"
)

type Permission string

const ReadSecret Permission = "read:secret"
const SetSecret Permission = "set:secret"
const DeleteSecret Permission = "delete:secret"
const DestroySecret Permission = "destroy:secret"

const ReadKey Permission = "read:key"
const SetKey Permission = "set:key"
const DeleteKey Permission = "delete:key"
const DestroyKey Permission = "destroy:key"
const SignKey Permission = "sign:key"
const EncryptKey Permission = "encrypt:key"

const ReadEth1 Permission = "read:eth1"
const SetEth1 Permission = "set:eth1"
const DeleteEth1 Permission = "delete:eth1"
const DestroyEth1 Permission = "destroy:eth1"
const SignEth1 Permission = "sign:eth1"
const EncryptEth1 Permission = "encrypt:eth1"

func ListPermissions() []Permission {
	return []Permission{
		ReadSecret,
		SetSecret,
		DeleteSecret,
		DestroySecret,
		ReadKey,
		SetKey,
		DeleteKey,
		DestroyKey,
		SignKey,
		EncryptKey,
		ReadEth1,
		SetEth1,
		DeleteEth1,
		DestroyEth1,
		SignEth1,
		EncryptEth1,
	}
}

func ListWildcardPermission(p string) []Permission {
	all := ListPermissions()
	parts := strings.Split(p, ":")
	action, resource := parts[0], parts[1]
	if action == "*" && resource == "*" {
		return all
	}

	included := []Permission{}
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
