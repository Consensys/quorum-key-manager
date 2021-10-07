package reader

import (
	"github.com/consensys/quorum-key-manager/src/infra/manifests"
	"github.com/consensys/quorum-key-manager/src/infra/manifests/entities"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type FilesystemReader struct {
	path  string
	isDir bool
}

var _ manifests.Reader = &FilesystemReader{}

func New(cfg *Config) (*FilesystemReader, error) {
	fs, err := os.Stat(cfg.Path)
	if err != nil {
		return nil, err
	}

	return &FilesystemReader{
		path:  cfg.Path,
		isDir: fs.IsDir(),
	}, nil
}

func (ll *FilesystemReader) Load() ([]*manifest.Manifest, error) {
	if !ll.isDir {
		return ll.buildMessages(ll.path)
	}

	var mnfs []*manifest.Manifest
	err := filepath.Walk(ll.path, func(fp string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileExtension := filepath.Ext(fp)
		if fileExtension == ".yml" || fileExtension == ".yaml" {
			currManifests, err := ll.buildMessages(fp)
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

func (ll *FilesystemReader) buildMessages(fp string) ([]*manifest.Manifest, error) {
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

	for _, mnf := range mnfs {
		err = val.Struct(mnf)
		if err != nil {
			return nil, err
		}
	}

	return mnfs, nil
}
