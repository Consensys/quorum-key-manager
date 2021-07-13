package interceptor

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/jsonrpc"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func (i *Interceptor) ethAccounts(ctx context.Context) ([]ethcommon.Address, error) {
	i.logger.Debug("listing ETH accounts")

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	storeAccounts, err := i.stores.ListAllAccounts(ctx, userInfo)
	if err != nil {
		return nil, err
	}

	addresses := []ethcommon.Address{}
	for _, storeAccount := range storeAccounts {
		addresses = append(addresses, storeAccount.Address)
	}

	i.logger.Debug("ETH accounts fetched successfully")
	return addresses, nil
}

func (i *Interceptor) EthAccounts() jsonrpc.Handler {
	h, _ := jsonrpc.MakeHandler(i.ethAccounts)
	return h
}
