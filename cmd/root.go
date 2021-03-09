package cmd

import (
	"os"

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

	return runCmd
}

func run(cmd *cobra.Command, _ []string) error {
	return nil
}
