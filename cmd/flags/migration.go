package flags

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(MigrationSourceURLViperKey, migrationSourceURLDefault)
	_ = viper.BindEnv(MigrationSourceURLViperKey, migrationSourceURLENV)
}

const (
	migrationSourceURLFlag     = "migration-source-url"
	MigrationSourceURLViperKey = "migration.source.url"
	migrationSourceURLDefault  = "/migrations"
	migrationSourceURLENV      = "MIGRATION_SOURCE_URL"
)

// MigrationFlags register flags for Postgres migrations
func MigrationFlags(f *pflag.FlagSet) {
	migrationSouceURL(f)
}

func migrationSouceURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Absolute path of the migrations directory
Environment variable: %q`, migrationSourceURLENV)
	f.String(migrationSourceURLFlag, migrationSourceURLDefault, desc)
	_ = viper.BindPFlag(MigrationSourceURLViperKey, f.Lookup(migrationSourceURLFlag))
}

func NewMigrationsConfig(vipr *viper.Viper) string {
	return vipr.GetString(MigrationSourceURLViperKey)
}
