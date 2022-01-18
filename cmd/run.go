package cmd

import (
	"os"

	"github.com/consensys/quorum-key-manager/cmd/flags"
	"github.com/consensys/quorum-key-manager/pkg/common"
	app "github.com/consensys/quorum-key-manager/src"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runCmd(cmd, args)
			if err != nil {
				cmd.SilenceUsage = true
			}
			return err
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			preRunBindFlags(viper.GetViper(), cmd.Flags(), "key-manager")
		},
	}

	flags.HTTPFlags(runCmd.Flags())
	flags.ManifestFlags(runCmd.Flags())
	flags.LoggerFlags(runCmd.Flags())
	flags.PGFlags(runCmd.Flags())
	flags.OIDCFlags(runCmd.Flags())
	flags.APIKeyFlags(runCmd.Flags())
	flags.TLSFlags(runCmd.Flags())

	return runCmd
}

func runCmd(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	vipr := viper.GetViper()
	cfg, err := flags.NewAppConfig(vipr)
	if err != nil {
		return err
	}

	logger, err := zap.NewLogger(cfg.Logger)
	if err != nil {
		return err
	}
	defer syncZapLogger(logger)

	appli, err := app.New(ctx, cfg, logger)
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
