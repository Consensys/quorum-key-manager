package hashicorp

import (
	hashicorp "github.com/hashicorp/vault/api"
	"io/ioutil"
	"strings"
)

func ManageToken(client *hashicorp.Client, cfg *Config) error {
	err := client.SetAddress(cfg.Address)
	if err != nil {
		return err
	}

	client.SetNamespace(cfg.Namespace)

	var decodedToken string
	if cfg.TokenFilePath != "" {
		encodedToken, err := ioutil.ReadFile(cfg.TokenFilePath)
		if err != nil {
			return err
		}

		decodedToken = strings.TrimSuffix(string(encodedToken), "\n") // Remove the newline if it exists
		decodedToken = strings.TrimSuffix(decodedToken, "\r")         // This one is for windows compatibility
	} else {
		decodedToken = cfg.Token
	}

	client.SetToken(decodedToken)

	return nil
}
