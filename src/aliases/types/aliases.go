package types

import (
	aliasentities "github.com/consensys/quorum-key-manager/src/aliases/entities"
)

type UpdateAliasRequest struct {
	aliasentities.Alias
}

type UpdateAliasResponse struct {
	aliasentities.Alias
}
