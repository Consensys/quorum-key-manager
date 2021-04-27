package flags

import (
	app "github.com/ConsenSysQuorum/quorum-key-manager/src"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
	"github.com/spf13/viper"
)

func NewAppConfig(vipr *viper.Viper) *app.Config {
	return &app.Config{
		Logger: newLoggerConfig(vipr),
		HTTP:   newHTTPConfig(vipr),
		Manifests: []*manifest.Manifest{
			newHashicorpManifest(vipr),
			newNodeManifest(vipr),
			newAKVManifest(vipr),
		},
	}
}
