package cmd

import (
	"fmt"
	"os"

	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/cmd/flags"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //nolint
	_ "github.com/golang-migrate/migrate/v4/source/file"       //nolint
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newMigrateCommand() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run migration",
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

	// Register Up command
	upCmd := &cobra.Command{
		Use:   "up [target]",
		Short: "Executes all migrations",
		Long:  "Executes all available migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrateUp()
		},
	}
	migrateCmd.AddCommand(upCmd)

	// Register Down command
	downCmd := &cobra.Command{
		Use:   "down",
		Short: "Reverts last migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrateDown()
		},
	}
	migrateCmd.AddCommand(downCmd)

	// Register Reset command
	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reverts all migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrateReset()
		},
	}
	migrateCmd.AddCommand(resetCmd)

	return migrateCmd
}

func migrateUp() error {
	vipr := viper.GetViper()
	logger, err := initLogger(vipr)
	if err != nil {
		return err
	}

	m, err := initMigrations(vipr, logger)
	if err != nil {
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

func migrateDown() error {
	vipr := viper.GetViper()
	logger, err := initLogger(vipr)
	if err != nil {
		return err
	}

	m, err := initMigrations(vipr, logger)
	if err != nil {
		return err
	}

	err = m.Steps(-1)
	if err != nil {
		errMessage := "failed to downgrade migrations"
		logger.WithError(err).Error(errMessage)
		return err
	}

	logger.Info("migration executed successfully")
	return nil
}

func migrateReset() error {
	vipr := viper.GetViper()
	logger, err := initLogger(vipr)
	if err != nil {
		return err
	}

	m, err := initMigrations(vipr, logger)
	if err != nil {
		return err
	}

	err = m.Down()
	if err != nil {
		errMessage := "failed to reset all migrations"
		logger.WithError(err).Error(errMessage)
		return err
	}

	logger.Info("migration executed successfully")
	return nil
}

func initMigrations(vipr *viper.Viper, logger log.Logger) (*migrate.Migrate, error) {
	pgCfg := flags.NewPostgresConfig(vipr)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", pgCfg.User, pgCfg.Password, pgCfg.Host, pgCfg.Port, pgCfg.Database)
	m, err := migrate.New("file:///migrations", dbURL)
	if err != nil {
		errMessage := "failed to create migration instance"
		logger.WithError(err).Error(errMessage)
		return nil, err
	}

	return m, nil
}

func initLogger(vipr *viper.Viper) (log.Logger, error) {
	logCfg := flags.NewLoggerConfig(vipr)
	return zap.NewLogger(logCfg)
}
