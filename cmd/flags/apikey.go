package flags

import (
	"fmt"

	"github.com/consensys/quorum-key-manager/src/infra/api-key/filesystem"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(authAPIKeyFileViperKey, authAPIKeyFileEnv)
}

const (
	authAPIKeyFileFlag        = "auth-api-key-file"
	authAPIKeyFileViperKey    = "auth.api.key.file"
	authAPIKeyDefaultFileFlag = ""
	authAPIKeyFileEnv         = "AUTH_API_KEY_FILE"
)

func APIKeyFlags(f *pflag.FlagSet) {
	authAPIKeyFile(f)
}

func authAPIKeyFile(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`API key CSV file location.
Environment variable: %q`, authAPIKeyFileEnv)
	f.String(authAPIKeyFileFlag, authAPIKeyDefaultFileFlag, desc)
	_ = viper.BindPFlag(authAPIKeyFileViperKey, f.Lookup(authAPIKeyFileFlag))
}

func NewAPIKeyConfig(vipr *viper.Viper) *filesystem.Config {
	path := vipr.GetString(authAPIKeyFileViperKey)

	if path != "" {
		return filesystem.NewConfig(path)
	}

	return nil
}
