package flags

import (
	app "github.com/ConsenSysQuorum/quorum-key-manager/src/app"
	"github.com/spf13/viper"
)

func NewAppConfig(vipr *viper.Viper) *app.Config {
	return &app.Config{
		Logger:       newLoggerConfig(vipr),
		HTTP:         newHTTPConfig(vipr),
		ManifestPath: newManifest(vipr),
	}
}
