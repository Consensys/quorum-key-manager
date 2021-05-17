package interceptor

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (i *Interceptor) ethSign(ctx context.Context, from ethcommon.Address, data hexutil.Bytes) (*hexutil.Bytes, error) {
	store, err := i.stores.GetAccountStoreByAddr(ctx, from)
	if err != nil {
		return nil, err
	}

	var sig hexutil.Bytes
	sig, err = store.Sign(ctx, from, data)
	if err != nil {
		return nil, err
	}

	return &sig, nil
}

func (i *Interceptor) EthSign() jsonrpc.Handler {
	h, _ := jsonrpc.MakeHandler(i.ethSign)
	return h
}
