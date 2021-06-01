package ethereum

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func init() {
	err := jsonrpc.ProvideCaller(
		privSrv,
	)
	if err != nil {
		panic(err)
	}
}

var privSrv = new(privService)

type privService struct {
	DistributeRawTransaction func(jsonrpc.Client) func(context.Context, hexutil.Bytes) (*hexutil.Bytes, error)                        `namespace:"priv"`
	GetEeaTransactionCount   func(jsonrpc.Client) func(context.Context, ethcommon.Address, string, []string) (*hexutil.Uint64, error) `namespace:"priv"`
	GetTransactionCount      func(jsonrpc.Client) func(context.Context, ethcommon.Address, string) (*hexutil.Uint64, error)           `namespace:"priv"`
}

//go:generate mockgen -source=caller_priv.go -destination=mock/caller_priv.go -package=mock

type PrivCaller interface {
	DistributeRawTransaction(context.Context, []byte) ([]byte, error)
	GetTransactionCount(ctx context.Context, addr ethcommon.Address, privacyGroupID string) (uint64, error)
	GetEeaTransactionCount(ctx context.Context, addr ethcommon.Address, privateFrom string, privateFor []string) (uint64, error)
}

type privCaller struct {
	client jsonrpc.Client
}

func (c *privCaller) DistributeRawTransaction(ctx context.Context, raw []byte) ([]byte, error) {
	enclaveKey, err := privSrv.DistributeRawTransaction(c.client)(ctx, hexutil.Bytes(raw))
	if err != nil {
		return nil, err
	}

	return []byte(*enclaveKey), nil
}

func (c *privCaller) GetTransactionCount(ctx context.Context, addr ethcommon.Address, privacyGroupID string) (uint64, error) {
	n, err := privSrv.GetTransactionCount(c.client)(ctx, addr, privacyGroupID)
	if err != nil {
		return 0, err
	}

	return uint64(*n), nil
}

func (c *privCaller) GetEeaTransactionCount(ctx context.Context, addr ethcommon.Address, privateFrom string, privateFor []string) (uint64, error) {
	n, err := privSrv.GetEeaTransactionCount(c.client)(ctx, addr, privateFrom, privateFor)
	if err != nil {
		return 0, err
	}

	return uint64(*n), nil
}
