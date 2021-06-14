package interceptor

import (
	"context"

	"github.com/consensysquorum/quorum-key-manager/pkg/jsonrpc"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (i *Interceptor) ethSign(ctx context.Context, from ethcommon.Address, data hexutil.Bytes) (*hexutil.Bytes, error) {
	store, err := i.stores.GetEth1StoreByAddr(ctx, from)
	if err != nil {
		return nil, err
	}

	sig, err := store.Sign(ctx, from.Hex(), data)
	if err != nil {
		return nil, err
	}

	return (*hexutil.Bytes)(&sig), nil
}

func (i *Interceptor) EthSign() jsonrpc.Handler {
	h, _ := jsonrpc.MakeHandler(i.ethSign)
	return h
}
