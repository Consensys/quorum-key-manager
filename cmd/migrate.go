package cmd

import (
	"fmt"
	"os"

	"github.com/consensys/quorum-key-manager/cmd/flags"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //nolint
	_ "github.com/golang-migrate/migrate/v4/source/file" //nolint
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newMigrateCommand() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run migration",
		RunE:  migrateUp,
		PreRun: func(cmd *cobra.Command, args []string) {
			preRunBindFlags(viper.GetViper(), cmd.Flags(), "key-manager")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			// TODO: Identify which error code to return
			os.Exit(0)
		},
	}

	flags.LoggerFlags(migrateCmd.Flags())
	flags.PGFlags(migrateCmd.Flags())

	return migrateCmd
}

func migrateUp(_ *cobra.Command, _ []string) error {
	vipr := viper.GetViper()
	pgCfg := flags.NewPostgresConfig(vipr)
	logCfg := flags.NewLoggerConfig(vipr)

	logger, err := zap.NewLogger(logCfg)
	if err != nil {
		return err
	}
	defer syncZapLogger(logger)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", pgCfg.User, pgCfg.Password, pgCfg.Host, pgCfg.Port, pgCfg.Database)
	m, err := migrate.New("file:///migrations", dbURL)
	if err != nil {
		errMessage := "failed to create migration instance"
		logger.WithError(err).Error(errMessage)
		return err
	}

	err = m.Up()
	if err != nil {
		errMessage := "failed to execute migrations"
		logger.WithError(err).Error(errMessage)
		return err
	}

	logger.Info("migration executed successfully")
	return nil
}
