package cmd

import (
	"fmt"

	"github.com/consensys/quorum-key-manager/cmd/flags"
	"github.com/consensys/quorum-key-manager/cmd/imports"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"
	manifestreader "github.com/consensys/quorum-key-manager/src/infra/manifests/filesystem"
	"github.com/consensys/quorum-key-manager/src/infra/postgres/client"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/database/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newImportCmd() *cobra.Command {
	var db database.Database
	var logger log.Logger
	var mnf *manifest.Manifest

	importCmd := &cobra.Command{
		Use:   "import",
		Short: "Import management tool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			logger, err = getLogger()
			if err != nil {
				return err
			}

			db, err = getDatabase(logger)
			if err != nil {
				return err
			}

			mnf, err = getManifest()
			if err != nil {
				return err
			}

			return nil
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			syncZapLogger(logger.(*zap.Logger))
		},
	}

	flags.PGFlags(importCmd.Flags())
	flags.ImportFlags(importCmd.Flags())
	flags.ManifestFlags(importCmd.Flags())

	importSecretsCmd := &cobra.Command{
		Use:   "secrets",
		Short: "import secrets from a vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			return imports.ImportSecrets(cmd.Context(), db.Secrets(mnf.Name), mnf, logger)
		},
	}
	importCmd.AddCommand(importSecretsCmd)

	importKeysCmd := &cobra.Command{
		Use:   "keys",
		Short: "import keys from a vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			return imports.ImportKeys(cmd.Context(), db.Keys(mnf.Name), mnf, logger)
		},
	}
	importCmd.AddCommand(importKeysCmd)

	importEthereumCmd := &cobra.Command{
		Use:   "ethereum",
		Short: "import ethereum accounts from a vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			return imports.ImportEthereum(cmd.Context(), db.ETHAccounts(mnf.Name), mnf, logger)
		},
	}
	importCmd.AddCommand(importEthereumCmd)

	return importCmd
}

func getLogger() (log.Logger, error) {
	return zap.NewLogger(flags.NewLoggerConfig(viper.GetViper()))
}

func getDatabase(logger log.Logger) (database.Database, error) {
	pgCfg, err := flags.NewPostgresConfig(viper.GetViper())
	if err != nil {
		return nil, err
	}

	// Create Postgres DB
	postgresClient, err := client.New(pgCfg)
	if err != nil {
		return nil, err
	}

	return postgres.New(logger, postgresClient), nil
}

func getManifest() (*manifest.Manifest, error) {
	vipr := viper.GetViper()
	// Get manifests
	manifestReader, err := manifestreader.New(flags.NewManifestConfig(vipr))
	if err != nil {
		return nil, err
	}

	manifests, err := manifestReader.Load()
	if err != nil {
		return nil, err
	}

	storeName := flags.GetStoreName(vipr)

	for _, mnf := range manifests {
		if mnf.Name == storeName {
			return mnf, nil
		}
	}

	return nil, fmt.Errorf("inexistent store %s in the manifests definitions", storeName)
}
