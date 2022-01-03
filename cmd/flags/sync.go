package flags

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(storeNameViperKey, storeNameDefault)
	_ = viper.BindEnv(storeNameViperKey, storeNameEnv)
}

const (
	storeNameFlag     = "sync-store-name"
	storeNameViperKey = "sync.store.name"
	storeNameDefault  = ""
	storeNameEnv      = "SYNC_STORE_NAME"
)

func SyncFlags(f *pflag.FlagSet) {
	storeName(f)
}

func storeName(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Name of the store to index
Environment variable: %q`, storeNameEnv)
	f.String(storeNameFlag, storeNameDefault, desc)
	_ = viper.BindPFlag(storeNameViperKey, f.Lookup(storeNameFlag))
}

func GetStoreName(vipr *viper.Viper) string {
	return vipr.GetString(storeNameViperKey)
}
