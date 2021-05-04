package interceptor

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func (i *Interceptor) ethAccounts(ctx context.Context) ([]ethcommon.Address, error) {
	storeAccounts, err := i.stores.ListAllAccounts(ctx)
	if err != nil {
		return nil, err
	}

	addresses := []ethcommon.Address{}
	for _, storeAccount := range storeAccounts {
		addresses = append(addresses, storeAccount.Address)
	}

	return addresses, nil
}

func (i *Interceptor) EthAccounts() jsonrpc.Handler {
	h, _ := jsonrpc.MakeHandler(i.ethAccounts)
	return h
}
