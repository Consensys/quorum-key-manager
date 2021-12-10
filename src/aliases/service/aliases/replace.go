package aliases

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (s *Aliases) Replace(ctx context.Context, addrs []string) ([]string, error) {
	var values []string
	for _, addr := range addrs {
		regName, aliasKey, isAlias := s.Parse(addr)

		// it is not an alias
		if !isAlias {
			values = append(values, addr)
			continue
		}

		alias, err := s.db.FindOne(ctx, regName, aliasKey)
		if err != nil {
			return nil, err
		}

		switch alias.Kind {
		case entities.AliasKindArray:
			vals, ok := alias.Value.([]interface{})
			if !ok {
				return nil, errors.InvalidFormatError("bad array format")
			}

			for _, v := range vals {
				str, ok := v.(string)
				if !ok {
					return nil, errors.InvalidFormatError("bad array value type")
				}

				values = append(values, str)
			}
		case entities.AliasKindString:
			values = append(values, alias.Value.(string))
		default:
			return nil, errors.InvalidFormatError("bad value kind")
		}

	}
	return values, nil
}

func (s *Aliases) ReplaceSimple(ctx context.Context, addr string) (string, error) {
	alias, err := s.Replace(ctx, []string{addr})
	if err != nil {
		return "", err
	}

	if len(alias) != 1 {
		s.logger.WithError(err).Error("wrong alias type")
		return "", errors.EncodingError("alias should only have 1 value")
	}

	return alias[0], nil
}
