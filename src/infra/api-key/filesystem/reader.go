package filesystem

import (
	"context"
	csv2 "encoding/csv"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/infra/api-key"
	"io"
	"os"
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
	fs os.FileInfo
}

var _ apikey.Reader = &Reader{}

func New(cfg *Config) (*Reader, error) {
	fs, err := os.Stat(cfg.Path)
	if err != nil {
		return nil, errors.ConfigError(err.Error())
	}

	return &Reader{fs: fs}, nil
}

func (r *Reader) Load(_ context.Context) (map[string]*entities.UserClaims, error) {
	csvfile, err := os.Open(r.fs.Name())
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
			return nil, errors.ConfigError(err.Error())
		}
		if len(cells) != csvRowLen {
			return nil, errors.ConfigError("invalid number of cells, should be %d", csvRowLen)
		}

		claims[cells[csvHashOffset]] = &entities.UserClaims{
			Subject: cells[csvUserOffset],
			Scope:   cells[csvPermissionsOffset],
			Roles:   cells[csvRolesOffset],
		}
	}

	return claims, nil
}
