package cmd

import (
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/jwt"
	"github.com/consensys/quorum-key-manager/pkg/tls/certificate"

	"github.com/consensys/quorum-key-manager/cmd/flags"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	sub        string
	scope      []string
	roles      []string
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
	flags.AuthOIDCClaimPermissions(generateJWTCmd.Flags())
	flags.AuthOIDCClaimRoles(generateJWTCmd.Flags())

	generateJWTCmd.Flags().StringVar(&sub, flags.AuthOIDCClaimUsernameDefault, "", "username and tenant added in claims")
	generateJWTCmd.Flags().StringArrayVar(&scope, flags.AuthOIDCClaimPermissionsDefault, []string{}, "permissions added in claims")
	generateJWTCmd.Flags().StringArrayVar(&roles, flags.AuthOIDCClaimRolesDefault, []string{}, "roles added in claims")
	generateJWTCmd.Flags().DurationVar(&expiration, "expiration", time.Hour, "token expiration time")

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

	keyFile := vipr.GetString(flags.AuthOIDCPrivKeyViperKey)

	_, err = os.Stat(keyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("cannot read OIDC Key file %s", keyFile)
		}
		return err
	}

	keyFileContent, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return err
	}

	oidcCfg := authCfg.OIDC

	privPem, _ := pem.Decode(keyFileContent)
	privPemBytes := privPem.Bytes

	signingKey, err := certificate.ParsePrivateKey(privPemBytes)
	if err != nil {
		return err
	}
	generator, err := jwt.NewTokenGenerator(signingKey)

	if err != nil {
		logger.Error("failed to generate access token", "err", err.Error())
		return err
	}

	token, err := generator.GenerateAccessToken(map[string]interface{}{
		oidcCfg.Claims.Subject: sub,
		oidcCfg.Claims.Scope:   strings.Join(scope, " "),
		oidcCfg.Claims.Roles:   strings.Join(roles, ","),
	}, expiration)
	if err != nil {
		logger.Error("failed to generate access token", "err", err.Error())
		return err
	}

	logger.Info("token generated", "value", token)
	return nil
}
