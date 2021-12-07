package yaml

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/consensys/quorum-key-manager/src/entities"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"

	"github.com/consensys/quorum-key-manager/src/infra/manifests"
)

type Reader struct {
	path     string
	isDir    bool
	validate *validator.Validate
}

var _ manifests.Reader = &Reader{}

func New(cfg *Config) (*Reader, error) {
	fs, err := os.Stat(cfg.Path)
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	err = validate.RegisterValidation("isManifestKind", isManifestKind)
	if err != nil {
		return nil, err
	}

	return &Reader{path: cfg.Path, isDir: fs.IsDir(), validate: validate}, nil
}

func (r *Reader) Load(_ context.Context) (map[string][]entities.Manifest, error) {
	manifestsMap := make(map[string][]entities.Manifest)

	if !r.isDir {
		mnfs, err := r.loadFile(r.path)
		if err != nil {
			return nil, err
		}

		addManifests(mnfs, manifestsMap)
		return manifestsMap, nil
	}

	err := filepath.Walk(r.path, func(fp string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileExtension := filepath.Ext(fp)
		if fileExtension == ".yml" || fileExtension == ".yaml" {
			mnfs, err := r.loadFile(fp)
			if err != nil {
				return err
			}

			addManifests(mnfs, manifestsMap)
			return nil
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return manifestsMap, nil
}

func (r *Reader) loadFile(fp string) ([]entities.Manifest, error) {
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	var mnfs []entities.Manifest
	err = yaml.Unmarshal(data, &mnfs)
	if err != nil {
		return nil, err
	}

	for _, mnf := range mnfs {
		err = r.validate.Struct(mnf)
		if err != nil {
			return nil, err
		}
	}

	return mnfs, nil
}

func addManifests(mnfs []entities.Manifest, manifestsMap map[string][]entities.Manifest) {
	for _, mnf := range mnfs {
		if _, ok := manifestsMap[mnf.Kind]; !ok {
			manifestsMap[mnf.Kind] = []entities.Manifest{mnf}
			continue
		}

		manifestsMap[mnf.Kind] = append(manifestsMap[mnf.Kind], mnf)
	}
}

func isManifestKind(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		switch fl.Field().String() {
		case entities.RoleKind, entities.StoreKind, entities.NodeKind, entities.VaultKind:
			return true
		default:
			return false
		}
	}

	return true
}
