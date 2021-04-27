package flags

import (
	"encoding/json"
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager/akv"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/types"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(akvEnvironmentViperKey, akvEnvironmentDefault)
	_ = viper.BindEnv(akvEnvironmentViperKey, AKVEnvironmentEnv)
}

const (
	AKVEnvironmentEnv      = "AKV_ENVIRONMENT"
	akvEnvironmentViperKey = "akv.environment"
	akvEnvironmentFlag     = "akv-environment"
	akvEnvironmentDefault  = "{}"
)

// Flags register flags for AKV
func AKVFlags(f *pflag.FlagSet) {
	akvEnvironment(f)
}

func akvEnvironment(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Specifies the AKV environment.
Environment variable: %q `, AKVEnvironmentEnv)
	f.String(akvEnvironmentFlag, akvEnvironmentDefault, desc)
	_ = viper.BindPFlag(akvEnvironmentViperKey, f.Lookup(akvEnvironmentFlag))
}

func newAKVManifest(vipr *viper.Viper) *manifest.Manifest {
	envStr := vipr.GetString(akvEnvironmentViperKey)
	specs := akv.Specs{}
	//TODO: Handle invalid not empty specs
	_ = json.Unmarshal([]byte(envStr), &specs)
	specRaw, _ := json.Marshal(specs)

	return &manifest.Manifest{
		Kind:    types.AKVSecrets,
		Name:    "AKVSecrets",
		Version: "0.0.1",
		Specs:   specRaw,
	}
}
