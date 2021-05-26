package cmd

import (
	"os"

	"github.com/ConsenSysQuorum/quorum-key-manager/cmd/flags"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	app "github.com/ConsenSysQuorum/quorum-key-manager/src/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCommand create root command
func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "quorum-kms",
		Short: "Run quorum-kms",
	}

	rootCmd.AddCommand(newRunCommand())

	return rootCmd
}

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		RunE:  run,
		PreRun: func(cmd *cobra.Command, args []string) {
			PreRunBindFlags(viper.GetViper(), cmd.Flags(), "quorum-kms")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			// TODO: Identify which error code to return
			os.Exit(0)
		},
	}

	flags.HTTPFlags(runCmd.Flags())
	flags.ManifestFlags(runCmd.Flags())
	flags.LoggerFlags(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, _ []string) error {
	vipr := viper.GetViper()
	cfg := flags.NewAppConfig(vipr)
	logger := log.NewLogger(cfg.Logger)
	log.SetDefaultLogger(logger)

	ctx := log.With(cmd.Context(), logger)
	appli, err := app.New(cfg, logger)
	if err != nil {
		logger.WithError(err).Error("application failed to start")
		return err
	}

	sig := common.NewSignalListener(func(sig os.Signal) {
		logger.WithField("sig", sig.String()).Warn("signal intercepted")
		if err = appli.Stop(ctx); err != nil {
			logger.WithError(err).Error("application stopped with errors")
		}
	})

	defer sig.Close()

	err = appli.Start(ctx)
	if err != nil {
		logger.WithError(err).Error("application exited with errors")
		return err
	}

	return nil
}
