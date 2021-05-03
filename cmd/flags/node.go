package flags

import (
	"encoding/json"
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
	proxynode "github.com/ConsenSysQuorum/quorum-key-manager/src/node/proxy"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(nodeAddrViperKey, nodeAddrDefault)
	_ = viper.BindEnv(nodeAddrViperKey, nodeAddrEnv)
}

const (
	NodeAddrFlag     = "node-addr"
	nodeAddrEnv      = "NODE_ADDR"
	nodeAddrViperKey = "node.addr"
	nodeAddrDefault  = "http://localhost:8545"
)

func nodeAddr(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Address of the JSON-RPC downstream node
Environment variable: %q`, nodeAddrEnv)
	f.String(NodeAddrFlag, nodeAddrDefault, desc)
	_ = viper.BindPFlag(nodeAddrViperKey, f.Lookup(NodeAddrFlag))
}

// NodeFlags register flags for Node
func NodeFlags(f *pflag.FlagSet) {
	nodeAddr(f)
}

func newNodeManifest(vipr *viper.Viper) *manifest.Manifest {
	specs := &proxynode.Config{
		RPC: &proxynode.DownstreamConfig{
			Addr: vipr.GetString(nodeAddrViperKey),
		},
	}

	specRaw, _ := json.Marshal(specs)

	return &manifest.Manifest{
		Kind:    "Node",
		Name:    "default",
		Version: "0.0.0",
		Specs:   specRaw,
	}
}
