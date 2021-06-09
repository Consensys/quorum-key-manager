package token

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/hashicorp"
	"github.com/fsnotify/fsnotify"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/hashicorp/vault/api"
)

// renewTokenLoop handle the token tokenWatcher of the application
type RenewTokenWatcher struct {
	tokenPath string
	client    hashicorp.VaultClient
	watcher   *fsnotify.Watcher
	logger    *log.Logger
}

func NewRenewTokenWatcher(client hashicorp.VaultClient, tokenPath string, logger *log.Logger) (*RenewTokenWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Add(filepath.Dir(tokenPath))
	if err != nil {
		return nil, errors.InvalidParameterError("cannot load token file at %s. %s", tokenPath, err.Error())
	}

	return &RenewTokenWatcher{
		tokenPath: tokenPath,
		client:    client,
		watcher:   watcher,
		logger:    logger.WithField("token_path", tokenPath),
	}, nil
}

// Run contains the token regeneration routine
func (rtl *RenewTokenWatcher) Start(ctx context.Context) error {
	defer rtl.watcher.Close()

	if err := rtl.refreshToken(); err != nil {
		return err
	}

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
				rtl.logger.WithField("event_name", event.Name).Debug("file has been updated")
				if err := rtl.refreshToken(); err != nil {
					return err
				}
			} else if event.Op&fsnotify.Create == fsnotify.Create {
				rtl.logger.WithField("event_name", event.Name).Debug("file has been created")
				if err := rtl.refreshToken(); err != nil {
					return err
				}
			}
			rtl.logger.Debug("event:", event)
		case err, ok := <-rtl.watcher.Errors:
			if !ok {
				return nil
			}
			rtl.logger.WithError(err).Error("failed to watch file events")
			return err
		}
	}
}

// Refresh the token
func (rtl *RenewTokenWatcher) refreshToken() error {
	encoded, err := ioutil.ReadFile(rtl.tokenPath)
	if err != nil {
		return errors.ConfigError("token file path could not be found at %s", rtl.tokenPath)
	}

	var wrappedToken api.SecretWrapInfo
	var token string
	err = json.Unmarshal(encoded, &wrappedToken)
	if err != nil {
		// Plain text token
		decoded := strings.TrimSuffix(string(encoded), "\n") // Remove the newline if it exists
		token = strings.TrimSuffix(decoded, "\r")            // This one is for windows compatibility
	} else {
		// Unwrap token
		secret, err2 := rtl.client.Client().Logical().Unwrap(wrappedToken.Token)
		if err2 != nil {
			return errors.InternalError("could not unwrap token")
		}
		token = fmt.Sprintf("%v", secret.Data["token"])
	}

	rtl.client.SetToken(token)
	rtl.logger.Info("token has been renewed")

	// Immediately delete the file after it was read
	err = os.Remove(rtl.tokenPath)
	if err != nil {
		rtl.logger.WithError(err).Warn("could not delete token file")
	}

	return nil
}
