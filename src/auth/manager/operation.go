package manager

type OpAction string
type OpResource string

var ActionRead OpAction = "read"
var ActionWrite OpAction = "write"
var ActionSign OpAction = "sign"
var ActionDelete OpAction = "delete"
var ActionDestroy OpAction = "destroy"

var ResourceKey OpResource = "key"
var ResourceSecret OpResource = "secret"
var ResourceEth OpResource = "eth1"
var ResourceNode OpResource = "node"

type Operation struct {
	Action   OpAction
	Resource OpResource
}
