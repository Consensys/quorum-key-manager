package cmd

import (
	"context"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/entities"
	storesservice "github.com/consensys/quorum-key-manager/src/stores"
	manifeststores "github.com/consensys/quorum-key-manager/src/stores/api/manifest"
	vaultsservice "github.com/consensys/quorum-key-manager/src/vaults"
	manifestvaults "github.com/consensys/quorum-key-manager/src/vaults/api/manifest"
	"github.com/consensys/quorum-key-manager/src/vaults/service/vaults"

	"github.com/consensys/quorum-key-manager/cmd/flags"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	manifestreader "github.com/consensys/quorum-key-manager/src/infra/manifests/yaml"
	"github.com/consensys/quorum-key-manager/src/infra/postgres/client"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newImportCmd() *cobra.Command {
	var logger *zap.Logger
	var vaultsService vaultsservice.Vaults
	var storesService storesservice.Stores
	var mnfs map[string][]entities.Manifest
	var storeName string

	userInfo := auth.NewWildcardUser()

	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Resource synchronization management tool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			storeName = flags.GetStoreName(viper.GetViper())

			var err error

			// Infra dependencies
			if logger, err = getLogger(); err != nil {
				return err
			}

			postgresClient, err := client.New(flags.NewPostgresConfig(viper.GetViper()))
			if err != nil {
				return err
			}

			if mnfs, err = getManifests(ctx); err != nil {
				return err
			}

			// Instantiate services
			vaultService := vaults.New(nil, logger)
			storesService = stores.NewConnector(nil, postgres.New(logger, postgresClient), vaultsService, logger)

			// Register vaults and stores
			if err = manifestvaults.NewVaultsHandler(vaultService).Register(ctx, mnfs[entities.VaultKind]); err != nil {
				return err
			}
			if err = manifeststores.NewStoresHandler(storesService).Register(ctx, mnfs[entities.StoreKind]); err != nil {
				return err
			}

			return nil
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			syncZapLogger(logger)
		},
	}

	flags.PGFlags(syncCmd.Flags())
	flags.SyncFlags(syncCmd.Flags())
	flags.ManifestFlags(syncCmd.Flags())

	syncSecretsCmd := &cobra.Command{
		Use:   "secrets",
		Short: "indexing secrets from remote vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			return storesService.ImportSecrets(cmd.Context(), storeName, userInfo)
		},
	}
	syncCmd.AddCommand(syncSecretsCmd)

	syncKeysCmd := &cobra.Command{
		Use:   "keys",
		Short: "indexing keys from remote vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			return storesService.ImportKeys(cmd.Context(), storeName, userInfo)
		},
	}
	syncCmd.AddCommand(syncKeysCmd)

	syncEthereumCmd := &cobra.Command{
		Use:   "ethereum",
		Short: "indexing ethereum accounts remote vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			return storesService.ImportKeys(cmd.Context(), storeName, userInfo)
		},
	}
	syncCmd.AddCommand(syncEthereumCmd)

	return syncCmd
}

func getLogger() (*zap.Logger, error) {
	return zap.NewLogger(flags.NewLoggerConfig(viper.GetViper()))
}

func getManifests(ctx context.Context) (map[string][]entities.Manifest, error) {
	manifestReader, err := manifestreader.New(flags.NewManifestConfig(viper.GetViper()))
	if err != nil {
		return nil, err
	}

	return manifestReader.Load(ctx)
}
