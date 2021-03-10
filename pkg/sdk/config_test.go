// +build unit

package client

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestKeyManagerTarget(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(flgs)
	expected := urlDefault
	assert.Equal(t, expected, viper.GetString(URLViperKey), "Default")

	_ = os.Setenv(urlEnv, "env-key-manager")
	expected = "env-key-manager"
	assert.Equal(t, expected, viper.GetString(URLViperKey), "From Environment Variable")
	_ = os.Unsetenv(urlEnv)

	args := []string{
		"--key-manager-url=flag-key-manager",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "Parse Key Manager flags should not error")
	expected = "flag-key-manager"
	assert.Equal(t, expected, viper.GetString(URLViperKey), "From Flag")
}

func TestFlags(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(f)
	assert.Equal(t, urlDefault, viper.GetString(URLViperKey), "Default")
}
