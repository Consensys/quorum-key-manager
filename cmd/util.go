package cmd

import (
	"os"
	"time"

	"github.com/consensys/quorum-key-manager/cmd/flags"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	"github.com/consensys/quorum-key-manager/pkg/tls/certificate"
	testutils2 "github.com/consensys/quorum-key-manager/pkg/tls/testutils"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/oicd"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/oicd/testutils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newUtilCommand() *cobra.Command {
	utilCmd := &cobra.Command{
		Use:   "util",
		Short: "Run util script",
		PostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	// Register Init command
	generateJWTCmd := &cobra.Command{
		Use:   "generate-token",
		Short: "Generate JWT Access Token",
		RunE:  runGenerateJWT,
	}

	utilCmd.AddCommand(generateJWTCmd)

	return utilCmd
}

func runGenerateJWT(_ *cobra.Command, _ []string) error {
	vipr := viper.GetViper()
	cfg := flags.NewAppConfig(vipr)

	logger, err := zap.NewLogger(cfg.Logger)
	if err != nil {
		return err
	}
	defer syncZapLogger(logger)

	oicdCfg := oicd.NewDefaultConfig()
	generator, err := testutils.NewJWTGenerator(&certificate.KeyPair{
		Cert: []byte(oicdCfg.Certificate),
		Key:  []byte(testutils2.OneLineRSAKeyPEMA),
	}, oicdCfg.Claims)

	if err != nil {
		logger.Error("failed to generate access token", "err", err.Error())
		return err
	}
	
	token, err := generator.GenerateAccessToken("username", []string{"group-admin", "no-existing-group"}, time.Minute * 60)
	if err != nil {
		logger.Error("failed to generate access token", "err", err.Error())
		return err
	}
	
	logger.Info("token generated", "value", token)
	return nil
}
