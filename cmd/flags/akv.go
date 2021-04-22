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
	_ = viper.BindEnv(akvEnvironmentViperKey, akvEnvironmentEnv)

	viper.SetDefault(akvResourceViperKey, akvResourceDefault)
	_ = viper.BindEnv(akvResourceViperKey, akvResourceEnv)
}

const (
	akvEnvironmentEnv      = "AZURE_ENVIRONMENT"
	akvEnvironmentViperKey = "azure.environment"
	akvEnvironmentFlag     = "hashicorp-token-file"
	akvEnvironmentDefault  = ""
)

const (
	akvResourceEnv      = "AZURE_KEYVAULT_RESOURCE"
	akvResourceFlag     = "azure-keyvault-resource"
	akvResourceViperKey = "azure.keyvault.resource"
	akvResourceDefault  = ""
)

// Flags register flags for AKV
func AKVFlags(f *pflag.FlagSet) {
	akvEnvironment(f)
	akvResource(f)
}

func akvEnvironment(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Specifies the AKV environment.
Environment variable: %q `, akvEnvironmentEnv)
	f.String(akvEnvironmentFlag, akvEnvironmentDefault, desc)
	_ = viper.BindPFlag(akvEnvironmentViperKey, f.Lookup(akvEnvironmentFlag))
}

func akvResource(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Specifies AKV resource.
Environment variable: %q `, akvResourceEnv)
	f.String(akvResourceFlag, akvResourceDefault, desc)
	_ = viper.BindPFlag(akvResourceViperKey, f.Lookup(akvResourceFlag))
}

func newAKVManifest(vipr *viper.Viper) *manifest.Manifest {
	specs := akv.Specs{
		EnvironmentName: vipr.GetString(akvEnvironmentViperKey),
		Resource:        vipr.GetString(akvResourceViperKey),
	}

	specRaw, _ := json.Marshal(specs)

	return &manifest.Manifest{
		Kind:    types.AKVSecrets,
		Name:    "AKVSecrets",
		Version: "0.0.0",
		Specs:   specRaw,
	}
}
