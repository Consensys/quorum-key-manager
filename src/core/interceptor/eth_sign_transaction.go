package interceptor

import (
	"context"
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	proxynode "github.com/ConsenSysQuorum/quorum-key-manager/src/node/proxy"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (i *Interceptor) ethSignTransaction(ctx context.Context, msg *ethereum.SendTxMsg) (*hexutil.Bytes, error) {
	if msg.Gas == nil {
		return nil, jsonrpc.InvalidParamsError(fmt.Errorf("gas not specified"))
	}

	if msg.GasPrice == nil {
		return nil, jsonrpc.InvalidParamsError(fmt.Errorf("gasPrice not specified"))
	}

	if msg.Nonce == nil {
		return nil, jsonrpc.InvalidParamsError(fmt.Errorf("nonce not specified"))
	}

	// Get store for from
	store, err := i.stores.GetEth1StoreByAddr(ctx, msg.From)
	if err != nil {
		return nil, err
	}

	// Get ChainID from Node
	sess := proxynode.SessionFromContext(ctx)
	chainID, err := sess.EthCaller().Eth().ChainID(ctx)
	if err != nil {
		return nil, err
	}

	// Sign
	var sig string
	if msg.IsPrivate() {
		sig, err = store.SignPrivate(ctx, msg.From.Hex(), msg.TxData())
	} else {
		sig, err = store.SignTransaction(ctx, msg.From.Hex(), chainID.String(), msg.TxData())
	}
	if err != nil {
		return nil, err
	}

	sigB, err := hexutil.Decode(sig)
	if err != nil {
		return nil, err
	}

	return (*hexutil.Bytes)(&sigB), nil
}

func (i *Interceptor) EthSignTransaction() jsonrpc.Handler {
	h, _ := jsonrpc.MakeHandler(i.ethSignTransaction)
	return h
}
