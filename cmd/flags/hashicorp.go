package flags

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager/hashicorp"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/types"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(hashicorpTokenFilePathViperKey, hashicorpTokenFilePathDefault)
	_ = viper.BindEnv(hashicorpTokenFilePathViperKey, hashicorpTokenFilePathEnv)

	viper.SetDefault(hashicorpMountPointViperKey, hashicorpMountPointDefault)
	_ = viper.BindEnv(hashicorpMountPointViperKey, hashicorpMountPointEnv)

	viper.SetDefault(hashicorpAddrViperKey, hashicorpAddrDefault)
	_ = viper.BindEnv(hashicorpAddrViperKey, hashicorpAddrEnv)
}

const (
	hashicorpTokenFilePathEnv      = "HASHICORP_TOKEN_FILE"
	hashicorpTokenFilePathDefault  = "/hashicorp/token/.hashicorp-token"
	hashicorpTokenFilePathViperKey = "hashicorp.token.file"
	HashicorpTokenFilePathFlag     = "hashicorp-token-file"
)

const (
	hashicorpMountPointEnv      = "HASHICORP_MOUNT_POINT"
	hashicorpMountPointFlag     = "hashicorp-mount-point"
	hashicorpMountPointViperKey = "hashicorp.mount.point"
	hashicorpMountPointDefault  = "orchestrate"
)

const (
	HashicorpAddrFlag     = "hashicorp-addr"
	hashicorpAddrEnv      = "HASHICORP_ADDR"
	hashicorpAddrViperKey = "hashicorp.addr"
	hashicorpAddrDefault  = "https://127.0.0.1:8200"
)

// Flags register flags for HashiCorp Hashicorp
func HashicorpFlags(f *pflag.FlagSet) {
	hashicorpAddr(f)
	hashicorpMountPoint(f)
	hashicorpTokenFilePath(f)
}

func hashicorpTokenFilePath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Specifies the token file path.
Parameter ignored if the token has been passed by HASHICORP_TOKEN
Environment variable: %q `, hashicorpTokenFilePathEnv)
	f.String(HashicorpTokenFilePathFlag, hashicorpTokenFilePathDefault, desc)
	_ = viper.BindPFlag(hashicorpTokenFilePathViperKey, f.Lookup(HashicorpTokenFilePathFlag))
}

func hashicorpMountPoint(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Specifies the mount point used. Should not start with a //
Environment variable: %q `, hashicorpMountPointEnv)
	f.String(hashicorpMountPointFlag, hashicorpMountPointDefault, desc)
	_ = viper.BindPFlag(hashicorpMountPointViperKey, f.Lookup(hashicorpMountPointFlag))
}

func hashicorpAddr(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp URL of the remote hashicorp hashicorp
Environment variable: %q`, hashicorpAddrEnv)
	f.String(HashicorpAddrFlag, hashicorpAddrDefault, desc)
	_ = viper.BindPFlag(hashicorpAddrViperKey, f.Lookup(HashicorpAddrFlag))
}

// ConfigFromViper returns a local config object that be converted into an api.Config
func newHashicorpManifest(vipr *viper.Viper) *manifest.Manifest {
	specs := hashicorp.SecretSpecs{
		MountPoint: vipr.GetString(hashicorpMountPointViperKey),
		Address:    vipr.GetString(hashicorpAddrViperKey),
	}

	tokenFilePath := vipr.GetString(hashicorpTokenFilePathViperKey)
	token, err := ioutil.ReadFile(tokenFilePath)
	if err == nil {
		specs.Token = string(token)
	}

	specRaw, _ := json.Marshal(specs)

	return &manifest.Manifest{
		Kind:    types.HashicorpSecrets,
		Name:    "HashicorpSecrets",
		Version: "0.0.0",
		Specs:   specRaw,
	}
}
