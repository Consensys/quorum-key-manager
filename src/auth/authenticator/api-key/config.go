package apikey

import (
	"encoding/base64"
	"hash"
)

type Config struct {
	APIKeyFile map[string]UserClaims
	Hasher     *hash.Hash
	B64Encoder *base64.Encoding
}

type UserClaims struct {
	UserName string
	Claims   []string
	Roles    []string
}

func NewConfig(apiKeyFile map[string]UserClaims, b64Encoder *base64.Encoding, hasher hash.Hash) *Config {
	return &Config{
		APIKeyFile: apiKeyFile,
		Hasher:     &hasher,
		B64Encoder: b64Encoder,
	}
}
