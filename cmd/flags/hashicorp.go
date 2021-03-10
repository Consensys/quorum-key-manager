package flags

import (
	"fmt"
	"time"

	hashicorp "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(hashicorpTokenViperKey, hashicorpTokenDefault)
	_ = viper.BindEnv(hashicorpTokenViperKey, hashicorpTokenEnv)

	viper.SetDefault(hashicorpTokenFilePathViperKey, hashicorpTokenFilePathDefault)
	_ = viper.BindEnv(hashicorpTokenFilePathViperKey, hashicorpTokenFilePathEnv)

	viper.SetDefault(hashicorpMountPointViperKey, hashicorpMountPointDefault)
	_ = viper.BindEnv(hashicorpMountPointViperKey, hashicorpMountPointEnv)

	viper.SetDefault(hashicorpRateLimitViperKey, hashicorpRateLimitDefault)
	_ = viper.BindEnv(hashicorpRateLimitViperKey, hashicorpRateLimitEnv)

	viper.SetDefault(hashicorpBurstLimitViperKey, hashicorpBurstLimitDefault)
	_ = viper.BindEnv(hashicorpBurstLimitViperKey, hashicorpBurstLimitEnv)

	viper.SetDefault(hashicorpAddrViperKey, hashicorpAddrDefault)
	_ = viper.BindEnv(hashicorpAddrViperKey, hashicorpAddrEnv)

	viper.SetDefault(hashicorpCACertViperKey, hashicorpCACertDefault)
	_ = viper.BindEnv(hashicorpCACertViperKey, hashicorpCACertEnv)

	viper.SetDefault(hashicorpCAPathViperKey, hashicorpCAPathDefault)
	_ = viper.BindEnv(hashicorpCAPathViperKey, hashicorpCAPathEnv)

	viper.SetDefault(hashicorpClientCertViperKey, hashicorpClientCertDefault)
	_ = viper.BindEnv(hashicorpClientCertViperKey, hashicorpClientCertEnv)

	viper.SetDefault(hashicorpClientKeyViperKey, hashicorpClientKeyDefault)
	_ = viper.BindEnv(hashicorpClientKeyViperKey, hashicorpClientKeyEnv)

	viper.SetDefault(hashicorpClientTimeoutViperKey, hashicorpClientTimeoutDefault)
	_ = viper.BindEnv(hashicorpClientTimeoutViperKey, hashicorpClientTimeoutEnv)

	viper.SetDefault(hashicorpMaxRetriesViperKey, hashicorpMaxRetriesDefault)
	_ = viper.BindEnv(hashicorpMaxRetriesViperKey, hashicorpMaxRetriesEnv)

	viper.SetDefault(hashicorpSkipVerifyViperKey, hashicorpSkipVerifyDefault)
	_ = viper.BindEnv(hashicorpSkipVerifyViperKey, hashicorpSkipVerifyEnv)

	viper.SetDefault(hashicorpTLSServerNameViperKey, hashicorpTLSServerNameDefault)
	_ = viper.BindEnv(hashicorpTLSServerNameViperKey, hashicorpTLSServerNameEnv)
}

const (
	hashicorpTokenEnv         = "HASHICORP_TOKEN"
	hashicorpTokenFilePathEnv = "HASHICORP_TOKEN_FILE"
	hashicorpMountPointEnv    = "HASHICORP_MOUNT_POINT"
	hashicorpRateLimitEnv     = "HASHICORP_RATE_LIMIT"
	hashicorpBurstLimitEnv    = "HASHICORP_BURST_LIMIT"
	hashicorpAddrEnv          = "HASHICORP_ADDR"
	hashicorpCACertEnv        = "HASHICORP_CACERT"
	hashicorpCAPathEnv        = "HASHICORP_CAPATH"
	hashicorpClientCertEnv    = "HASHICORP_CLIENT_CERT"
	hashicorpClientKeyEnv     = "HASHICORP_CLIENT_KEY"
	hashicorpClientTimeoutEnv = "HASHICORP_CLIENT_TIMEOUT"
	hashicorpMaxRetriesEnv    = "HASHICORP_MAX_RETRIES"
	hashicorpSkipVerifyEnv    = "HASHICORP_SKIP_VERIFY"
	hashicorpTLSServerNameEnv = "HASHICORP_TLS_SERVER_NAME"

	HashicorpTokenFlag         = "hashicorp-token"
	HashicorpTokenFilePathFlag = "hashicorp-token-file"
	hashicorpMountPointFlag    = "hashicorp-mount-point"
	hashicorpRateLimitFlag     = "hashicorp-rate-limit"
	hashicorpBurstLimitFlag    = "hashicorp-burst-limit"
	HashicorpAddrFlag          = "hashicorp-addr"
	hashicorpCACertFlag        = "hashicorp-cacert"
	hashicorpCAPathFlag        = "hashicorp-capath"
	hashicorpClientCertFlag    = "hashicorp-client-cert"
	hashicorpClientKeyFlag     = "hashicorp-client-key"
	hashicorpClientTimeoutFlag = "hashicorp-client-timeout"
	hashicorpMaxRetriesFlag    = "hashicorp-max-retries"
	hashicorpSkipVerifyFlag    = "hashicorp-skip-verify"
	hashicorpTLSServerNameFlag = "hashicorp-tls-server-name"

	hashicorpTokenViperKey         = "hashicorp.token"
	hashicorpTokenFilePathViperKey = "hashicorp.token.file"
	hashicorpMountPointViperKey    = "hashicorp.mount.point"
	hashicorpRateLimitViperKey     = "hashicorp.rate.limit"
	hashicorpBurstLimitViperKey    = "hashicorp.burst.limit"
	hashicorpAddrViperKey          = "hashicorp.addr"
	hashicorpCACertViperKey        = "hashicorp.cacert"
	hashicorpCAPathViperKey        = "hashicorp.capath"
	hashicorpClientCertViperKey    = "hashicorp.client.cert"
	hashicorpClientKeyViperKey     = "hashicorp.client.key"
	hashicorpClientTimeoutViperKey = "hashicorp.client.timeout"
	hashicorpMaxRetriesViperKey    = "hashicorp.max.retries"
	hashicorpSkipVerifyViperKey    = "hashicorp.skip.verify"
	hashicorpTLSServerNameViperKey = "hashicorp.tls.server.name"

	// No need to redefine the default here
	hashicorpTokenDefault         = ""
	hashicorpTokenFilePathDefault = "/hashicorp/token/.hashicorp-token"
	hashicorpMountPointDefault    = "orchestrate"
	hashicorpRateLimitDefault     = float64(0)
	hashicorpBurstLimitDefault    = int(0)
	hashicorpAddrDefault          = "https://127.0.0.1:8200"
	hashicorpCACertDefault        = ""
	hashicorpCAPathDefault        = ""
	hashicorpClientCertDefault    = ""
	hashicorpClientKeyDefault     = ""
	hashicorpClientTimeoutDefault = time.Second * 60
	hashicorpMaxRetriesDefault    = int(0)
	hashicorpSkipVerifyDefault    = false
	hashicorpTLSServerNameDefault = ""
)

// Flags register flags for HashiCorp Hashicorp
func HashicorpFlags(f *pflag.FlagSet) {
	hashicorpAddr(f)
	hashicorpBurstLimit(f)
	hashicorpCACert(f)
	hashicorpCAPath(f)
	hashicorpClientCert(f)
	hashicorpClientKey(f)
	hashicorpClientTimeout(f)
	hashicorpMaxRetries(f)
	hashicorpMountPoint(f)
	hashicorpRateLimit(f)
	hashicorpSkipVerify(f)
	hashicorpTLSServerName(f)
	hashicorpTokenFilePath(f)
	hashicorpToken(f)
}

func hashicorpToken(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Specifies the static token. Parameter ignored if the token has been passed by HASHICORP_TOKEN.
Environment variable: %q `, hashicorpTokenEnv)
	f.String(HashicorpTokenFlag, hashicorpTokenDefault, desc)
	_ = viper.BindPFlag(hashicorpTokenViperKey, f.Lookup(HashicorpTokenFlag))
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

func hashicorpRateLimit(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp query rate limit
Environment variable: %q`, hashicorpRateLimitEnv)
	f.Float64(hashicorpRateLimitFlag, hashicorpRateLimitDefault, desc)
	_ = viper.BindPFlag(hashicorpRateLimitViperKey, f.Lookup(hashicorpRateLimitFlag))
}

func hashicorpBurstLimit(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp query burst limit
Environment variable: %q`, hashicorpRateLimitEnv)
	f.Int(hashicorpBurstLimitFlag, hashicorpBurstLimitDefault, desc)
	_ = viper.BindPFlag(hashicorpBurstLimitViperKey, f.Lookup(hashicorpBurstLimitFlag))
}

func hashicorpAddr(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp URL of the remote hashicorp hashicorp
Environment variable: %q`, hashicorpAddrEnv)
	f.String(HashicorpAddrFlag, hashicorpAddrDefault, desc)
	_ = viper.BindPFlag(hashicorpAddrViperKey, f.Lookup(HashicorpAddrFlag))
}

func hashicorpCACert(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp CA certificate
Environment variable: %q`, hashicorpCACertEnv)
	f.String(hashicorpCACertFlag, hashicorpCACertDefault, desc)
	_ = viper.BindPFlag(hashicorpCACertViperKey, f.Lookup(hashicorpCACertFlag))
}

func hashicorpCAPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Path toward the CA certificate
Environment variable: %q`, hashicorpCAPathEnv)
	f.String(hashicorpCAPathFlag, hashicorpCAPathDefault, desc)
	_ = viper.BindPFlag(hashicorpCAPathViperKey, f.Lookup(hashicorpCAPathFlag))
}

func hashicorpClientCert(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Certificate of the client
Environment variable: %q`, hashicorpClientCertEnv)
	f.String(hashicorpClientCertFlag, hashicorpClientCertDefault, desc)
	_ = viper.BindPFlag(hashicorpClientCertViperKey, f.Lookup(hashicorpClientCertFlag))
}

func hashicorpClientKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp client key
Environment variable: %q`, hashicorpClientKeyEnv)
	f.String(hashicorpClientKeyFlag, hashicorpClientKeyDefault, desc)
	_ = viper.BindPFlag(hashicorpClientKeyViperKey, f.Lookup(hashicorpClientKeyFlag))
}

func hashicorpClientTimeout(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp clean timeout of the client
Environment variable: %q`, hashicorpClientTimeoutEnv)
	f.Duration(hashicorpClientTimeoutFlag, hashicorpClientTimeoutDefault, desc)
	_ = viper.BindPFlag(hashicorpClientTimeoutViperKey, f.Lookup(hashicorpClientTimeoutFlag))
}

func hashicorpMaxRetries(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp max retry for a request
Environment variable: %q`, hashicorpMaxRetriesEnv)
	f.Int(hashicorpMaxRetriesFlag, hashicorpMaxRetriesDefault, desc)
	_ = viper.BindPFlag(hashicorpMaxRetriesViperKey, f.Lookup(hashicorpMaxRetriesFlag))
}

func hashicorpSkipVerify(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp skip verification
Environment variable: %q`, hashicorpSkipVerifyEnv)
	f.Bool(hashicorpSkipVerifyFlag, hashicorpSkipVerifyDefault, desc)
	_ = viper.BindPFlag(hashicorpSkipVerifyViperKey, f.Lookup(hashicorpSkipVerifyFlag))
}

func hashicorpTLSServerName(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp TLS server name
Environment variable: %q`, hashicorpTLSServerNameEnv)
	f.String(hashicorpTLSServerNameFlag, hashicorpTLSServerNameDefault, desc)
	_ = viper.BindPFlag(hashicorpTLSServerNameViperKey, f.Lookup(hashicorpTLSServerNameFlag))
}

// ConfigFromViper returns a local config object that be converted into an api.Config
func NewHashicorpConfig() *hashicorp.Config {
	return &hashicorp.Config{
		Address:       viper.GetString(hashicorpAddrViperKey),
		BurstLimit:    viper.GetInt(hashicorpBurstLimitViperKey),
		CACert:        viper.GetString(hashicorpCACertViperKey),
		CAPath:        viper.GetString(hashicorpCAPathViperKey),
		ClientCert:    viper.GetString(hashicorpClientCertViperKey),
		ClientKey:     viper.GetString(hashicorpClientKeyViperKey),
		ClientTimeout: viper.GetDuration(hashicorpClientTimeoutViperKey),
		MaxRetries:    viper.GetInt(hashicorpMaxRetriesViperKey),
		MountPoint:    viper.GetString(hashicorpMountPointViperKey),
		RateLimit:     viper.GetFloat64(hashicorpRateLimitViperKey),
		SkipVerify:    viper.GetBool(hashicorpSkipVerifyViperKey),
		TLSServerName: viper.GetString(hashicorpTLSServerNameViperKey),
		TokenFilePath: viper.GetString(hashicorpTokenFilePathViperKey),
		Token:         viper.GetString(hashicorpTokenViperKey),
		Renewable:     true,
	}
}
