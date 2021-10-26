package csv

import (
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/infra/api-key"
	"hash"
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
	userClaims map[string]*entities.UserClaims
	hasher     hash.Hash
}

var _ apikey.Reader = &Reader{}

func New(cfg *Config) (*Reader, error) {
	// Open the file
	csvfile, err := os.Open(cfg.Path)
	if err != nil {
		return nil, err
	}
	defer csvfile.Close()

	// Parse the file
	r := csv.NewReader(csvfile)
	// Set separator
	r.Comma = csvSeparator
	// ignore comments in file
	r.Comment = csvCommentsMarker

	retFile := make(map[string]*entities.UserClaims)

	// Iterate through the lines
	for {
		// Read each line from csv
		cells, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(cells) != csvRowLen {
			return nil, fmt.Errorf("invalid number of cells, should be %d", csvRowLen)
		}

		retFile[cells[csvHashOffset]] = &entities.UserClaims{
			Subject: cells[csvUserOffset],
			Scope:   cells[csvPermissionsOffset],
			Roles:   cells[csvRolesOffset],
		}
	}

	return &Reader{
		userClaims: retFile,
		hasher:     sha256.New(),
	}, nil
}

func (reader *Reader) Get(_ context.Context, apiKey []byte) (*entities.UserClaims, error) {
	h := reader.hasher
	h.Reset()
	_, err := h.Write(apiKey)
	if err != nil {
		return nil, errors.UnauthorizedError(err.Error())
	}
	clientAPIKeyHash := h.Sum(nil)

	strClientHash := hex.EncodeToString(clientAPIKeyHash)
	claims, ok := reader.userClaims[strClientHash]
	if !ok {
		return nil, errors.UnauthorizedError("invalid api-key")
	}

	return claims, nil
}
