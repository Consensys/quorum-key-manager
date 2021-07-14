package flags

import (
	"fmt"
	"os"

	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(manifestPathKey, manifestPathDefault)
	_ = viper.BindEnv(manifestPathKey, manifestPathEnv)
}

const (
	ManifestPath        = "manifest-path"
	manifestPathEnv     = "MANIFEST_PATH"
	manifestPathKey     = "manifest.path"
	manifestPathDefault = ""
)

func manifestPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Path to manifest file/folder to configure key manager stores and nodes
Environment variable: %q`, manifestPathEnv)
	f.String(ManifestPath, manifestPathDefault, desc)
	_ = viper.BindPFlag(manifestPathKey, f.Lookup(ManifestPath))
}

// ManifestFlags register flags for Node
func ManifestFlags(f *pflag.FlagSet) {
	manifestPath(f)
}

func newManifestsConfig(vipr *viper.Viper) (*manifestsmanager.Config, error) {
	manifestPath := vipr.GetString(manifestPathKey)
	_, err := os.Stat(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("cannot read manifest path %s", manifestPath)
		}
		return nil, err
	}

	return &manifestsmanager.Config{
		Path: vipr.GetString(manifestPathKey),
	}, nil
}
