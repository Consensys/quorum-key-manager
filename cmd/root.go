package cmd

import (
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/consensys/quorum-key-manager/src/infra/log/zap"

	"github.com/consensys/quorum-key-manager/cmd/flags"
	"github.com/consensys/quorum-key-manager/pkg/common"
	app "github.com/consensys/quorum-key-manager/src"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCommand create root command
func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "key-manager",
		Short: "Run Quorum Key Manager",
	}

	rootCmd.AddCommand(newRunCommand())
	rootCmd.AddCommand(newMigrateCmd())

	return rootCmd
}

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		RunE:  run,
		PreRun: func(cmd *cobra.Command, args []string) {
			PreRunBindFlags(viper.GetViper(), cmd.Flags(), "key-manager")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			// TODO: Identify which error code to return
			os.Exit(0)
		},
	}

	flags.HTTPFlags(runCmd.Flags())
	flags.ManifestFlags(runCmd.Flags())
	flags.LoggerFlags(runCmd.Flags())
	flags.PGFlags(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	vipr := viper.GetViper()
	cfg := flags.NewAppConfig(vipr)

	logger, err := zap.NewLogger(cfg.Logger)
	if err != nil {
		return err
	}
	defer syncZapLogger(logger)

	appli, err := app.New(cfg, logger)
	if err != nil {
		logger.WithError(err).Error("could not create app")
		return err
	}

	done := make(chan struct{})
	sig := common.NewSignalListener(func(sig os.Signal) {
		logger.With("sig", sig.String()).Warn("signal intercepted")
		if err = appli.Stop(ctx); err != nil {
			logger.WithError(err).Error("application stopped with errors")
		}
		close(done)
	})

	defer sig.Close()

	err = appli.Start(ctx)
	if err != nil {
		logger.WithError(err).Error("application failed to start")
		return err
	}

	<-done

	return nil
}

func syncZapLogger(logger *zap.Logger) {
	_ = logger.Sync()
}

// newMigrateCmd create migrate command
func newMigrateCmd() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrate()
		},
	}

	// Register Init command
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrate("init")
		},
	}
	migrateCmd.AddCommand(initCmd)

	// Register Up command
	upCmd := &cobra.Command{
		Use:   "up [target]",
		Short: "Upgrade database",
		Long:  "Runs all available migrations or up to [target] if argument is provided",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrate(append([]string{"up"}, args...)...)
		},
	}
	migrateCmd.AddCommand(upCmd)

	// Register Down command
	downCmd := &cobra.Command{
		Use:   "down",
		Short: "Reverts last migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrate("down")
		},
	}
	migrateCmd.AddCommand(downCmd)

	// Register Reset command
	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reverts all migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrate("reset")
		},
	}
	migrateCmd.AddCommand(resetCmd)

	return migrateCmd
}

func migrate(args ...string) error {
	// Set database connection
	opts, err := flags.NewPostgresConfig(viper.GetViper()).ToPGOptions()
	if err != nil {
		return err
	}
	db := pg.Connect(opts)

	oldVersion, newVersion, err := migrations.Run(db, args...)
	if err != nil {
		log.WithError(err).Errorf("Migration failed")
		return err
	}

	err = db.Close()
	if err != nil {
		log.WithError(err).Warn("could not close Postgres connection")
	}

	log.WithField("version", newVersion).WithField("previous_version", oldVersion).Info("All migrations completed")
	return nil
}
