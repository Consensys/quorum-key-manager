package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/consensys/quorum-key-manager/cmd/flags"
	"github.com/consensys/quorum-key-manager/pkg/auth"
	"github.com/consensys/quorum-key-manager/pkg/tls/certificate"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	username   string
	groups     []string
	expiration time.Duration
)

func newUtilCommand() *cobra.Command {
	utilCmd := &cobra.Command{
		Use:   "utils",
		Short: "Run util script",
		PostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	// Register Init command
	generateJWTCmd := &cobra.Command{
		Use:   "generate-jwt",
		Short: "Generate JWT Access Token",
		RunE:  runGenerateJWT,
	}

	flags.LoggerFlags(generateJWTCmd.Flags())
	flags.AuthOIDCClaimUsername(generateJWTCmd.Flags())
	flags.AuthOIDCClaimGroups(generateJWTCmd.Flags())
	flags.AuthOIDCCertKeyFile(generateJWTCmd.Flags())

	generateJWTCmd.Flags().StringVar(&username, "username", "", "username added in claims")
	generateJWTCmd.Flags().StringArrayVar(&groups, "groups", []string{}, "groups added in claims")
	generateJWTCmd.Flags().DurationVar(&expiration, "expiration", time.Hour, "Token expiration time")

	utilCmd.AddCommand(generateJWTCmd)

	return utilCmd
}

func runGenerateJWT(_ *cobra.Command, _ []string) error {
	vipr := viper.GetViper()
	authCfg, err := flags.NewAuthConfig(vipr)
	if err != nil {
		return err
	}

	loggerCfg := flags.NewLoggerConfig(vipr)
	logger, err := zap.NewLogger(loggerCfg)
	if err != nil {
		return err
	}
	defer syncZapLogger(logger)

	keyFile := vipr.GetString(flags.AuthOIDCCAKeyFileViperKey)
	_, err = os.Stat(keyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("cannot read CA Key file %s", keyFile)
		}
		return err
	}

	keyFileContent, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return err
	}

	oidcCfg := authCfg.OIDC
	var keys [][]byte
	keys, err = certificate.Decode(keyFileContent, "PRIVATE KEY")
	if err != nil {
		return err
	}

	certKey, err := certificate.ParsePrivateKey(keys[0])
	if err != nil {
		return err
	}

	generator, err := auth.NewJWTGenerator(certKey, oidcCfg.Claims)

	if err != nil {
		logger.Error("failed to generate access token", "err", err.Error())
		return err
	}

	token, err := generator.GenerateAccessToken(username, groups, expiration)
	if err != nil {
		logger.Error("failed to generate access token", "err", err.Error())
		return err
	}

	logger.Info("token generated", "value", token)
	return nil
}
