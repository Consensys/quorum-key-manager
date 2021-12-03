package cmd

import (
	"context"
	"fmt"
	"github.com/consensys/quorum-key-manager/src/entities"

	"github.com/consensys/quorum-key-manager/cmd/flags"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	manifestreader "github.com/consensys/quorum-key-manager/src/infra/manifests/yaml"
	"github.com/consensys/quorum-key-manager/src/infra/postgres/client"
	storeservice "github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newImportCmd() *cobra.Command {
	var logger *zap.Logger
	var storesConnector storeservice.Stores
	var mnf *entities.Manifest

	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Resource synchronization management tool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if logger, err = getLogger(); err != nil {
				return err
			}
			if storesConnector, err = getStores(logger); err != nil {
				return err
			}
			if mnf, err = getManifest(cmd.Context()); err != nil {
				return err
			}

			return nil
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			syncZapLogger(logger)
		},
	}

	flags.PGFlags(syncCmd.Flags())
	flags.ImportFlags(syncCmd.Flags())
	flags.ManifestFlags(syncCmd.Flags())

	syncSecretsCmd := &cobra.Command{
		Use:   "secrets",
		Short: "indexing secrets from remote vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if err := storesConnector.CreateSecret(ctx, mnf.Name, mnf.Kind, mnf.Specs, mnf.AllowedTenants); err != nil {
				return err
			}

			return storesConnector.ImportSecrets(cmd.Context(), mnf.Name, entities.NewWildcardUser())
		},
	}
	syncCmd.AddCommand(syncSecretsCmd)

	syncKeysCmd := &cobra.Command{
		Use:   "keys",
		Short: "indexing keys from remote vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if err := storesConnector.CreateKey(ctx, mnf.Name, entities2.VaultType(mnf.Kind), mnf.Specs, mnf.AllowedTenants); err != nil {
				return err
			}

			return storesConnector.ImportKeys(cmd.Context(), mnf.Name, entities.NewWildcardUser())
		},
	}
	syncCmd.AddCommand(syncKeysCmd)

	syncEthereumCmd := &cobra.Command{
		Use:   "ethereum",
		Short: "indexing ethereum accounts remote vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if err := storesConnector.CreateEthereum(ctx, mnf.Name, entities2.VaultType(mnf.Kind), mnf.Specs, mnf.AllowedTenants); err != nil {
				return err
			}

			return storesConnector.ImportEthereum(cmd.Context(), mnf.Name, entities.NewWildcardUser())
		},
	}
	syncCmd.AddCommand(syncEthereumCmd)

	return syncCmd
}

func getLogger() (*zap.Logger, error) {
	return zap.NewLogger(flags.NewLoggerConfig(viper.GetViper()))
}

func getStores(logger log.Logger) (storeservice.Stores, error) {
	// Create Postgres DB
	postgresClient, err := client.New(flags.NewPostgresConfig(viper.GetViper()))
	if err != nil {
		return nil, err
	}

	return stores.NewConnector(nil, postgres.New(logger, postgresClient), logger), nil
}

func getManifest(ctx context.Context) (*entities2.Manifest, error) {
	vipr := viper.GetViper()
	// Get manifests
	manifestReader, err := manifestreader.New(flags.NewManifestConfig(vipr))
	if err != nil {
		return nil, err
	}

	manifests, err := manifestReader.Load(ctx)
	if err != nil {
		return nil, err
	}

	storeName := flags.GetStoreName(vipr)

	for _, mnf := range manifests {
		// TODO: Filter on Load() function from reader when ManifestKind Store implemented
		if mnf.Kind == manifestreader.Role || mnf.Kind == manifestreader.Node {
			continue
		}

		if mnf.Name == storeName {
			return mnf, nil
		}
	}

	return nil, fmt.Errorf("inexistent store %s in the manifests definitions", storeName)
}
