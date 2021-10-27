package filesystem

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/consensys/quorum-key-manager/src/infra/manifests"
	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type Reader struct {
	fs os.FileInfo
}

var _ manifests.Reader = &Reader{}

func New(cfg *Config) (*Reader, error) {
	fs, err := os.Stat(cfg.Path)
	if err != nil {
		return nil, errors.ConfigError(err.Error())
	}

	return &Reader{fs: fs}, nil
}

func (r *Reader) Load(_ context.Context) ([]*manifest.Manifest, error) {
	if !r.fs.IsDir() {
		return r.buildMessages(r.fs.Name())
	}

	var mnfs []*manifest.Manifest
	err := filepath.Walk(r.fs.Name(), func(fp string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.ConfigError(err.Error())
		}

		if info.IsDir() {
			return nil
		}

		fileExtension := filepath.Ext(fp)
		if fileExtension == ".yml" || fileExtension == ".yaml" {
			currManifests, err := r.buildMessages(fp)
			if err != nil {
				return err
			}

			mnfs = append(mnfs, currManifests...)

			return nil
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return mnfs, nil
}

func (r *Reader) buildMessages(fp string) ([]*manifest.Manifest, error) {
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, errors.ConfigError(err.Error())
	}

	var mnfs []*manifest.Manifest
	err = yaml.Unmarshal(data, &mnfs)
	if err != nil {
		return nil, errors.InvalidFormatError(err.Error())
	}

	for _, mnf := range mnfs {
		err = validator.New().Struct(mnf)
		if err != nil {
			return nil, errors.InvalidFormatError(err.Error())
		}
	}

	return mnfs, nil
}
