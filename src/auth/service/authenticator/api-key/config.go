package apikey

import (
	"encoding/base64"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"hash"
)

type Config struct {
	APIKeyFile map[string]*entities.UserClaims
	Hasher     *hash.Hash
	B64Encoder *base64.Encoding
}

func NewConfig(apiKeyFile map[string]*entities.UserClaims, b64Encoder *base64.Encoding, hasher hash.Hash) *Config {
	return &Config{
		APIKeyFile: apiKeyFile,
		Hasher:     &hasher,
		B64Encoder: b64Encoder,
	}
}
