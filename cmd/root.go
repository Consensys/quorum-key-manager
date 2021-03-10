package cmd

import (
	"os"

	"github.com/ConsenSysQuorum/quorum-key-manager/cmd/flags"
	"github.com/ConsenSysQuorum/quorum-key-manager/src"
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

	flags.APIFlags(runCmd.Flags())
	flags.HashicorpFlags(runCmd.Flags())
	flags.LoggerFlags(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, _ []string) error {
	cfg := flags.NewAppConfig(viper.GetViper())
	app := src.New(cfg)
	return app.Start(cmd.Context())
}
