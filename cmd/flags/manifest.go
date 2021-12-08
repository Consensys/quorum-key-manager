package flags

import (
	"fmt"

	manifests "github.com/consensys/quorum-key-manager/src/infra/manifests/yaml"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(manifestPathViperKey, manifestPathDefault)
	_ = viper.BindEnv(manifestPathViperKey, manifestPathEnv)
}

const (
	ManifestPath         = "manifest-path"
	manifestPathEnv      = "MANIFEST_PATH"
	manifestPathViperKey = "manifest.path"
	manifestPathDefault  = ""
)

func manifestPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Path to manifest file/folder to configure key manager stores and nodes
Environment variable: %q`, manifestPathEnv)
	f.String(ManifestPath, manifestPathDefault, desc)
	_ = viper.BindPFlag(manifestPathViperKey, f.Lookup(ManifestPath))
}

// ManifestFlags register flags for Node
func ManifestFlags(f *pflag.FlagSet) {
	manifestPath(f)
}

func NewManifestConfig(vipr *viper.Viper) *manifests.Config {
	return manifests.NewConfig(vipr.GetString(manifestPathViperKey))
}
