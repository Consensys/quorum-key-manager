package manager

type OpAction string
type OpResource string

var ActionRead OpAction = "read"
var ActionWrite OpAction = "write"
var ActionSign OpAction = "sign"
var ActionDelete OpAction = "delete"
var ActionDestroy OpAction = "destroy"

var ResourceKey OpResource = "keys"
var ResourceSecret OpResource = "secrets"
var ResourceEth1Account OpResource = "eth1accounts"
var ResourceStores OpResource = "stores"
var ResourceNode OpResource = "nodes"

type Operation struct {
	Action   OpAction
	Resource OpResource
}
