package filesystem

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/consensys/quorum-key-manager/src/infra/manifests"
	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type Reader struct {
	path  string
	isDir bool
}

var _ manifests.Reader = &Reader{}

func New(cfg *Config) (*Reader, error) {
	fs, err := os.Stat(cfg.Path)
	if err != nil {
		return nil, err
	}

	return &Reader{path: cfg.Path, isDir: fs.IsDir()}, nil
}

func (r *Reader) Load(_ context.Context) ([]*manifest.Manifest, error) {
	if !r.isDir {
		return r.buildMessages(r.path)
	}

	var mnfs []*manifest.Manifest
	err := filepath.Walk(r.path, func(fp string, info os.FileInfo, err error) error {
		if err != nil {
			return err
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
		return nil, err
	}

	var mnfs []*manifest.Manifest
	err = yaml.Unmarshal(data, &mnfs)
	if err != nil {
		return nil, err
	}

	for _, mnf := range mnfs {
		err = validator.New().Struct(mnf)
		if err != nil {
			return nil, err
		}
	}

	return mnfs, nil
}
