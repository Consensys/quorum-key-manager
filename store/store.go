package store

import (
	"context"

	"github.com/ConsenSys/quorum-signer/store/types"
)

// Store is responsible to store items
type Store interface {
	// Info returns secret store information
	Info(context.Context) *types.StoreInfo
}
