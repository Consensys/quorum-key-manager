package cmd

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func PreRunBindFlags(vipr *viper.Viper, flgs *pflag.FlagSet, ignore string) {
	for _, vk := range vipr.AllKeys() {
		if ignore != "" && strings.HasPrefix(vk, ignore) {
			continue
		}

		// Convert viperKey to cmd flag name
		// For example: 'rest.api' to "rest-api"
		name := strings.Replace(vk, ".", "-", -1)

		// Only bind in case command flags contain the name
		if flgs.Lookup(name) != nil {
			_ = viper.BindPFlag(vk, flgs.Lookup(name))
		}
	}
}
