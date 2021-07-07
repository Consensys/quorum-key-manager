package flags

import (
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/spf13/viper"
)

func newAuthConfig(_ *viper.Viper) *auth.Config {
	return &auth.Config{}
}
