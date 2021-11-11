package csv

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	csv2 "encoding/csv"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/consensys/quorum-key-manager/src/auth/entities"
	apikey "github.com/consensys/quorum-key-manager/src/infra/api-key"
)

const (
	csvSeparator         = ','
	csvCommentsMarker    = '#'
	csvRowLen            = 4
	csvHashOffset        = 0
	csvUserOffset        = 1
	csvPermissionsOffset = 2
	csvRolesOffset       = 3
)

type Reader struct {
	path   string
	hasher hash.Hash
}

var _ apikey.Reader = &Reader{}

func New(cfg *Config) (*Reader, error) {
	_, err := os.Stat(cfg.Path)
	if err != nil {
		return nil, err
	}

	return &Reader{path: cfg.Path, hasher: sha256.New()}, nil
}

func (r *Reader) Load(_ context.Context) (map[string]*entities.UserClaims, error) {
	csvfile, err := os.Open(r.path)
	if err != nil {
		return nil, err
	}
	defer csvfile.Close()

	csvReader := csv2.NewReader(csvfile)
	csvReader.Comma = csvSeparator
	csvReader.Comment = csvCommentsMarker

	// Read each line from csv and fill claims
	claims := make(map[string]*entities.UserClaims)
	for {
		cells, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(cells) != csvRowLen {
			return nil, fmt.Errorf("invalid number of cells, should be %d", csvRowLen)
		}

		r.hasher.Reset()
		_, err = r.hasher.Write([]byte(cells[csvHashOffset]))
		if err != nil {
			return nil, fmt.Errorf("failed to hash api key")
		}

		apiKeyHash := base64.StdEncoding.EncodeToString(r.hasher.Sum(nil))
		claims[apiKeyHash] = &entities.UserClaims{
			Subject: cells[csvUserOffset],
			Scope:   cells[csvPermissionsOffset],
			Roles:   cells[csvRolesOffset],
		}
	}

	return claims, nil
}
