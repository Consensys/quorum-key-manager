package token

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/consensys/quorum-key-manager/src/infra/hashicorp"
	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/fsnotify/fsnotify"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/hashicorp/vault/api"
)

// RenewTokenWatcher handle the token tokenWatcher of the application
type RenewTokenWatcher struct {
	tokenPath     string
	client        hashicorp.Client
	watcher       *fsnotify.Watcher
	logger        log.Logger
	isTokenLoaded bool
}

func NewRenewTokenWatcher(client hashicorp.Client, tokenPath string, logger log.Logger) (*RenewTokenWatcher, error) {
	logger = logger.With("token_path", tokenPath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		errMessage := "failed to instantiate watcher"
		logger.WithError(err).Error(errMessage)
		return nil, errors.DependencyFailureError(errMessage)
	}

	err = watcher.Add(filepath.Dir(tokenPath))
	if err != nil {
		errMessage := "failed to load token file. Please verify the token file path"
		logger.WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	return &RenewTokenWatcher{
		tokenPath: tokenPath,
		client:    client,
		watcher:   watcher,
		logger:    logger,
	}, nil
}

// Start contains the token regeneration routine
func (rtl *RenewTokenWatcher) Start(ctx context.Context) error {
	defer rtl.watcher.Close()

	// First token refresh
	if err := rtl.refreshToken(); err != nil {
		return err
	}

	rtl.isTokenLoaded = true

	for {
		select {
		case <-ctx.Done():
			return nil
		case event, ok := <-rtl.watcher.Events:
			if !ok {
				return nil
			}

			if event.Name != rtl.tokenPath {
				continue
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				rtl.logger.Debug("token file has been updated")
				if err := rtl.refreshToken(); err != nil {
					return err
				}
			} else if event.Op&fsnotify.Create == fsnotify.Create {
				rtl.logger.Debug("token file has been created")
				if err := rtl.refreshToken(); err != nil {
					return err
				}
			}
		case err, ok := <-rtl.watcher.Errors:
			if !ok {
				return nil
			}
			rtl.logger.WithError(err).Error("failed to watch file events")
			return err
		}
	}
}

func (rtl *RenewTokenWatcher) IsTokenLoaded() bool {
	return rtl.isTokenLoaded
}

func (rtl *RenewTokenWatcher) refreshToken() error {
	encoded, err := ioutil.ReadFile(rtl.tokenPath)
	if err != nil {
		errMessage := "token file path could not be found"
		rtl.logger.WithError(err).Error(errMessage)
		return errors.ConfigError("token file path could not be found")
	}

	var wrappedToken api.SecretWrapInfo
	var token string
	err = json.Unmarshal(encoded, &wrappedToken)
	if err != nil {
		// Plain text token
		decoded := strings.TrimSuffix(string(encoded), "\n") // Delete the newline if it exists
		token = strings.TrimSuffix(decoded, "\r")            // This one is for windows compatibility
	} else {
		// Unwrap token
		secret, err2 := rtl.client.UnwrapToken(wrappedToken.Token)
		if err2 != nil {
			errMessage := "could not unwrap token"
			rtl.logger.WithError(err2).Error(errMessage)
			return errors.HashicorpVaultError(errMessage)
		}
		token = fmt.Sprintf("%v", secret.Data["token"])
	}

	rtl.client.SetToken(token)

	rtl.logger.Info("token has been successfully renewed")
	return nil
}
