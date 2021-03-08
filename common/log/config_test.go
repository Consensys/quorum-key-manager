// +build unit

package log

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLevel(t *testing.T) {
	name := LevelViperKey
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(flgs)

	expected := "info"
	if viper.GetString(name) != expected {
		t.Errorf("Level #1: expected %q but got %q", expected, viper.GetString(name))
	}

	_ = os.Setenv(levelEnv, "fatal")
	expected = "fatal"
	if viper.GetString(name) != expected {
		t.Errorf("Level #2: expected %q but got %q", expected, viper.GetString(name))
	}
	_ = os.Unsetenv(levelEnv)

	args := []string{
		fmt.Sprintf("--%s=%s", levelFlag, "debug"),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "debug"
	if viper.GetString(name) != expected {
		t.Errorf("Level #3: expected %q but got %q", expected, viper.GetString(name))
	}
}

func TestFormat(t *testing.T) {
	name := FormatViperKey
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(f)
	expected := "text"
	if viper.GetString(name) != expected {
		t.Errorf("Format #1: expected %q but got %q", expected, viper.GetString(name))
	}

	_ = os.Setenv(formatEnv, "json")
	expected = "json"
	if viper.GetString(name) != expected {
		t.Errorf("Format #2: expected %q but got %q", expected, viper.GetString(name))
	}
	_ = os.Unsetenv(formatEnv)

	args := []string{
		fmt.Sprintf("--%s=%s", formatFlag, "xml"),
	}
	err := f.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "xml"
	if viper.GetString(name) != expected {
		t.Errorf("Format #3: expected %q but got %q", expected, viper.GetString(name))
	}
}

func TestTimestamp(t *testing.T) {
	name := TimestampViperKey
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(f)
	if viper.GetBool(name) == true {
		t.Errorf("Timestamp #1: expected %v but got %v", false, viper.GetBool(name))
	}

	_ = os.Setenv(timestampEnv, "true")
	if viper.GetBool(name) != true {
		t.Errorf("Timestamp #2: expected %v but got %v", true, viper.GetBool(name))
	}
	_ = os.Unsetenv(timestampEnv)
	
	args := []string{
		fmt.Sprintf("--%s=%s", timestampFlag, "true"),
	}
	err := f.Parse(args)
	assert.NoError(t, err, "No error expected")
	
	if viper.GetBool(name) != true {
		t.Errorf("Timestamp #3: expected %v but got %v", true, viper.GetBool(name))
	}
}
