package cmd

import (
	"github.com/consensys/quorum-key-manager/cmd/flags"
	"github.com/consensys/quorum-key-manager/cmd/imports"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	"github.com/consensys/quorum-key-manager/src/infra/postgres/client"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/database/postgres"
	"github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newImportCmd() *cobra.Command {
	var db database.Database
	var logger log.Logger

	accountCmd := &cobra.Command{
		Use:   "import",
		Short: "Import management tool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			vipr := viper.GetViper()
			pgCfg, err := flags.NewPostgresConfig(vipr)
			if err != nil {
				return err
			}

			// Create logger
			logger, err = zap.NewLogger(flags.NewLoggerConfig(vipr))
			if err != nil {
				return err
			}

			// Create Postgres DB
			postgresClient, err := client.New(pgCfg)
			if err != nil {
				return err
			}

			db = postgres.New(logger, postgresClient)

			return nil
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			syncZapLogger(logger.(*zap.Logger))
		},
	}

	flags.PGFlags(accountCmd.Flags())
	flags.ImportFlags(accountCmd.Flags())

	importSecretsCmd := &cobra.Command{
		Use:   "secrets",
		Short: "import secrets from a vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			vipr := viper.GetViper()

			return imports.ImportSecrets(cmd.Context(), flags.GetStoreName(vipr), db, mock.NewMockSecretStore(nil))
		},
	}
	accountCmd.AddCommand(importSecretsCmd)

	importKeysCmd := &cobra.Command{
		Use:   "keys",
		Short: "import keys from a vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			vipr := viper.GetViper()
			return imports.ImportKeys(cmd.Context(), flags.GetStoreName(vipr), db, mock.NewMockKeyStore(nil))
		},
	}
	accountCmd.AddCommand(importKeysCmd)

	importEthereumCmd := &cobra.Command{
		Use:   "ethereum",
		Short: "import ethereum accounts from a vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			vipr := viper.GetViper()
			return imports.ImportEthereum(cmd.Context(), flags.GetStoreName(vipr), db, mock.NewMockEthStore(nil))
		},
	}
	accountCmd.AddCommand(importEthereumCmd)

	return accountCmd
}
