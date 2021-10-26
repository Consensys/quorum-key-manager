package flags

import (
	app "github.com/consensys/quorum-key-manager/src"
	"github.com/spf13/viper"
)

func NewAppConfig(vipr *viper.Viper) (*app.Config, error) {
	httpCfg, err := newHTTPConfig(vipr)
	if err != nil {
		return nil, err
	}

	return &app.Config{
		Logger:   NewLoggerConfig(vipr),
		HTTP:     httpCfg,
		Manifest: NewManifestConfig(vipr),
		OIDC:     NewOIDCConfig(vipr),
		APIKey:   NewAPIKeyConfig(vipr),
		Postgres: NewPostgresConfig(vipr),
	}, nil
}
