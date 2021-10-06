package manager

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/entities"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type LocalManager struct {
	path   string
	isDir  bool
	logger log.Logger
}

func NewLocalManager(path string, logger log.Logger) (*LocalManager, error) {
	fs, err := os.Stat(path)
	if err == nil {
		return &LocalManager{
			path:   path,
			isDir:  fs.IsDir(),
			logger: logger,
		}, nil
	}

	if os.IsNotExist(err) {
		errMessage := "folder or file does not exists"
		logger.WithError(err).Error(errMessage, "path", path)
		return nil, errors.InvalidParameterError(errMessage)
	}

	return nil, err
}

func (ll *LocalManager) Load() ([]manifest.Message, error) {
	logger := ll.logger.With("path", ll.path, "isDir", ll.isDir)
	logger.Debug("reading manifest items")

	if !ll.isDir {
		return ll.buildMessages(ll.path)
	}

	var messages []manifest.Message
	err := filepath.Walk(ll.path, func(fp string, info os.FileInfo, err error) error {
		if err != nil {
			errMessage := "failed to walk the file tree"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		fileExtension := filepath.Ext(fp)
		if fileExtension == ".yml" || fileExtension == ".yaml" {
			messages, err = ll.buildMessages(fp)
			if err != nil {
				errMessage := "failed to load manifests from file, please verify the manifest file format"
				logger.WithError(err).Error(errMessage)
				return errors.InvalidFormatError(errMessage)
			}

			return nil
		}

		errMessage := "unrecognised manifest extension, should be YAML"
		logger.Error(errMessage)
		return fmt.Errorf(errMessage)
	})
	if err != nil {
		errMessage := "failed to load manifests from file"
		logger.WithError(err).Error(errMessage)
		return nil, errors.InvalidFormatError(errMessage)
	}

	return messages, nil
}

func (ll *LocalManager) buildMessages(fp string) ([]manifest.Message, error) {
	val := validator.New()
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	var mnfs []*manifest.Manifest
	err = yaml.Unmarshal(data, &mnfs)
	if err != nil {
		return nil, err
	}

	var messages []manifest.Message
	for _, mnf := range mnfs {
		err = val.Struct(mnf)
		if err != nil {
			return nil, err
		}

		messages = append(messages, manifest.Message{
			Action:   manifest.CreateAction,
			Manifest: mnf,
		})
	}

	return messages, nil
}
