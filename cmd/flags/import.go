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
	storeNameFlag     = "import-store-name"
	storeNameViperKey = "import.store.name"
	storeNameDefault  = ""
	storeNameEnv      = "IMPORT_STORE_NAME"
)

func ImportFlags(f *pflag.FlagSet) {
	storeName(f)
}

func storeName(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Name of the store
Environment variable: %q`, storeNameEnv)
	f.String(storeNameFlag, storeNameDefault, desc)
	_ = viper.BindPFlag(storeNameViperKey, f.Lookup(storeNameFlag))
}

func GetStoreName(vipr *viper.Viper) string {
	return vipr.GetString(storeNameViperKey)
}
