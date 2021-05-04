package flags

import (
	"encoding/json"
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager/hashicorp"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/types"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(hashicorpEnvironmentViperKey, hashicorpEnvironmentDefault)
	_ = viper.BindEnv(hashicorpEnvironmentViperKey, hashicorpEnvironmentEnv)
}

const (
	hashicorpEnvironmentEnv      = "HASHICORP_ENVIRONMENT"
	hashicorpEnvironmentViperKey = "hashicorp.environment"
	hashicorpEnvironmentFlag     = "hashicorp-environment"
	hashicorpEnvironmentDefault  = "{}"
)

// Flags register flags for HashiCorp Hashicorp
func HashicorpFlags(f *pflag.FlagSet) {
	hashicorpEnvironment(f)
}

func hashicorpEnvironment(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Specifies the Hashicorp environment.
Environment variable: %q `, hashicorpEnvironmentEnv)
	f.String(hashicorpEnvironmentFlag, hashicorpEnvironmentDefault, desc)
	_ = viper.BindPFlag(hashicorpEnvironmentViperKey, f.Lookup(hashicorpEnvironmentFlag))
}

func newHashicorpSecretsManifest(vipr *viper.Viper) *manifest.Manifest {
	envStr := vipr.GetString(hashicorpEnvironmentViperKey)
	specs := hashicorp.SecretSpecs{}
	//TODO: Handle invalid not empty specs
	_ = json.Unmarshal([]byte(envStr), &specs)
	specRaw, _ := json.Marshal(specs)

	return &manifest.Manifest{
		Kind:    types.HashicorpSecrets,
		Name:    "HashicorpSecrets",
		Version: "0.0.1",
		Specs:   specRaw,
	}
}

func newHashicorpKeysManifest(vipr *viper.Viper) *manifest.Manifest {
	envStr := vipr.GetString(hashicorpEnvironmentViperKey)
	specs := hashicorp.SecretSpecs{}
	//TODO: Handle invalid not empty specs
	_ = json.Unmarshal([]byte(envStr), &specs)
	specRaw, _ := json.Marshal(specs)

	return &manifest.Manifest{
		Kind:    types.HashicorpKeys,
		Name:    "HashicorpKeys",
		Version: "0.0.1",
		Specs:   specRaw,
	}
}
