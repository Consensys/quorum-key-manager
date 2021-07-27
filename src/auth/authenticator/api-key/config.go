package apikey

import (
	"encoding/base64"
	"hash"
)

type Config struct {
	APIKeyFile map[string]UserNameAndGroups
	Hasher     *hash.Hash
	B64Encoder *base64.Encoding
}

type UserNameAndGroups struct {
	UserName string
	Groups   []string
}

func NewConfig(apiKeyFile map[string]UserNameAndGroups, b64Encoder *base64.Encoding, hasher hash.Hash) *Config {
	return &Config{
		APIKeyFile: apiKeyFile,
		Hasher:     &hasher,
		B64Encoder: b64Encoder,
	}
}
