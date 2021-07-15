package apikey

import "hash"

type Config struct {
	ApiKeyFile map[string]*UserNameAndGroups
	Hasher     hash.Hash
}

type UserNameAndGroups struct {
	UserName string
	Groups   []string
}

func NewConfig(apiKeyFile map[string]*UserNameAndGroups) *Config {
	return &Config{
		ApiKeyFile: apiKeyFile,
	}
}
