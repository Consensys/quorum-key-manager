package cmd

import (
	"os"

	"github.com/consensys/quorum-key-manager/cmd/flags"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	models2 "github.com/consensys/quorum-key-manager/src/stores/database/models"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newMigrateCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run migration",
		RunE:  migrateCmd,
		PreRun: func(cmd *cobra.Command, args []string) {
			preRunBindFlags(viper.GetViper(), cmd.Flags(), "key-manager")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			// TODO: Identify which error code to return
			os.Exit(0)
		},
	}

	flags.LoggerFlags(runCmd.Flags())
	flags.PGFlags(runCmd.Flags())

	return runCmd
}

func migrateCmd(cmd *cobra.Command, _ []string) error {
	vipr := viper.GetViper()
	pgCfg := flags.NewPostgresConfig(vipr)
	logCfg := flags.NewLoggerConfig(vipr)

	logger, err := zap.NewLogger(logCfg)
	if err != nil {
		return err
	}
	defer syncZapLogger(logger)

	pgOpt, err := pgCfg.ToPGOptions()
	if err != nil {
		return err
	}

	db := pg.Connect(pgOpt)
	defer db.Close()

	opts := &orm.CreateTableOptions{
		FKConstraints: true,
	}
	// we create tables for each model
	for _, v := range []interface{}{
		&models2.Secret{},
		&models2.Key{},
		&models2.ETH1Account{},
		&aliasent.Alias{},
	} {
		err = db.Model(v).CreateTable(opts)
		if err != nil {
			return err
		}
	}

	logger.Info("migration executed successfully")
	return nil
}
