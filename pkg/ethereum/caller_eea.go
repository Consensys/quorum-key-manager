package ethereum

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func init() {
	err := jsonrpc.ProvideCaller(
		eeaSrv,
	)
	if err != nil {
		panic(err)
	}
}

var eeaSrv = new(eeaService)

// eeaService is a jsonrpc.Caller which methods are meant to be automatically populed using jsonrpc.ProvideCaller
type eeaService struct {
	SendRawTransaction func(jsonrpc.Client) func(context.Context, hexutil.Bytes) (ethcommon.Hash, error) `namespace:"eea"`
}

//go:generate mockgen -source=caller_eea.go -destination=mock/caller_eea.go -package=mock

// EEACaller is a JSON-RPC client to a Ethereum client using eea namespace
type EEACaller interface {
	SendRawTransaction(context.Context, []byte) (ethcommon.Hash, error)
}

type eeaCaller struct {
	client jsonrpc.Client
}

func (c *eeaCaller) SendRawTransaction(ctx context.Context, raw []byte) (ethcommon.Hash, error) {
	return eeaSrv.SendRawTransaction(c.client)(ctx, hexutil.Bytes(raw))
}
