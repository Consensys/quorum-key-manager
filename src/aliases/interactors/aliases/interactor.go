package interactors

import (
	"context"
	"regexp"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/aliases"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log"
)

// We make sure Connector implements aliases.Interactor
var _ aliases.Interactor = &Interactor{}

// Interactor is the service layer for other service to query.
type Interactor struct {
	db    aliases.Interactor
	regex *regexp.Regexp

	logger log.Logger
}

func NewInteractor(db aliases.Interactor, logger log.Logger) (*Interactor, error) {
	const aliasParseFormat = `{{(?m)(?P<registry>[a-zA-Z0-9-_+]+):(?P<alias>[a-zA-Z0-9-_+]+)}}$`
	regex, err := regexp.Compile(aliasParseFormat)
	if err != nil {
		return nil, errors.ConfigError("bad regexp format '%v': %v", aliasParseFormat, err)
	}
	return &Interactor{
		db:    db,
		regex: regex,

		logger: logger,
	}, nil
}

func (i *Interactor) CreateAlias(ctx context.Context, registry string, alias aliasent.Alias) (*aliasent.Alias, error) {
	logger := i.logger.With(
		"registry_name", registry,
		"alias_key", alias.Key,
	)
	a, err := i.db.CreateAlias(ctx, registry, alias)
	if err != nil {
		return nil, err
	}
	logger.Info("alias created successfully")
	return a, nil
}

func (i *Interactor) GetAlias(ctx context.Context, registry, aliasKey string) (*aliasent.Alias, error) {
	return i.db.GetAlias(ctx, registry, aliasKey)
}

func (i *Interactor) UpdateAlias(ctx context.Context, registry string, alias aliasent.Alias) (*aliasent.Alias, error) {
	logger := i.logger.With(
		"registry_name", registry,
		"alias_key", alias.Key,
	)
	a, err := i.db.UpdateAlias(ctx, registry, alias)
	if err != nil {
		return nil, err
	}
	logger.Info("alias updated successfully")
	return a, nil
}

func (i *Interactor) DeleteAlias(ctx context.Context, registry, aliasKey string) error {
	logger := i.logger.With(
		"registry_name", registry,
		"alias_key", aliasKey,
	)
	err := i.db.DeleteAlias(ctx, registry, aliasKey)
	if err != nil {
		return err
	}
	logger.Info("alias deleted successfully")
	return nil
}

func (i *Interactor) ListAliases(ctx context.Context, registry string) ([]aliasent.Alias, error) {
	return i.db.ListAliases(ctx, registry)
}

func (i *Interactor) DeleteRegistry(ctx context.Context, registry string) error {
	logger := i.logger.With(
		"registry_name", registry,
	)
	err := i.db.DeleteRegistry(ctx, registry)
	if err != nil {
		return err
	}
	logger.Info("registry deleted successfully")
	return nil
}

// ParseAlias parses an alias string and returns the registryName and the aliasKey
// as well as if the string isAlias. If the string is not isAlias, we'll consider it
// as a valid key.
func (i *Interactor) ParseAlias(alias string) (regName, aliasKey string, isAlias bool) {
	submatches := i.regex.FindStringSubmatch(alias)
	if len(submatches) < 3 {
		return "", "", false
	}

	regName = submatches[1]
	aliasKey = submatches[2]

	return regName, aliasKey, true
}

// ReplaceAliases replace a slice of potential aliases with a slice having all the aliases replaced by their value.
// It will fail if no aliases can be found.
func (i *Interactor) ReplaceAliases(ctx context.Context, addrs []string) ([]string, error) {
	var values []string
	for _, addr := range addrs {
		regName, aliasKey, isAlias := i.ParseAlias(addr)

		// it is not an alias
		if !isAlias {
			values = append(values, addr)
			continue
		}

		alias, err := i.db.GetAlias(ctx, regName, aliasKey)
		if err != nil {
			return nil, err
		}

		values = append(values, alias.Value...)
	}
	return values, nil
}
