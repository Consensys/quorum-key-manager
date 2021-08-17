package flags

import (
	app "github.com/consensys/quorum-key-manager/src"
	"github.com/spf13/viper"
)

func NewAppConfig(vipr *viper.Viper) (*app.Config, error) {
	authCfg, err := NewAuthConfig(vipr)
	if err != nil {
		return nil, err
	}

	manifestCfg, err := newManifestsConfig(vipr)
	if err != nil {
		return nil, err
	}
	
	httpCfg, err := newHTTPConfig(vipr)
	if err != nil {
		return nil, err
	}

	return &app.Config{
		Logger:    NewLoggerConfig(vipr),
		HTTP:      httpCfg,
		Manifests: manifestCfg,
		Auth:      authCfg,
		Postgres:  NewPostgresConfig(vipr),
	}, nil
}
