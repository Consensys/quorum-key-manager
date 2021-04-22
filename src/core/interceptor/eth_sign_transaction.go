package interceptor

import (
	"context"
	"math/big"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/node"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (i *Interceptor) ethSignTransaction(ctx context.Context, msg *ethereum.SendTxMsg) (*hexutil.Bytes, error) {
	// Get store for from
	store, err := i.stores.GetAccountStoreByAddr(ctx, msg.From)
	if err != nil {
		return nil, err
	}

	txData, err := msg.TxData()
	if err != nil {
		return nil, err
	}

	// Get ChainID from Node
	sess := node.SessionFromContext(ctx)
	chainID, err := sess.EthClient().Eth().ChainID(ctx)
	if err != nil {
		return nil, err
	}

	// Sign
	var sig hexutil.Bytes
	if msg.IsPrivate() {
		sig, err = store.SignPrivate(ctx, msg.From, txData)
	} else {
		sig, err = store.SignEIP155(ctx, (*big.Int)(chainID), msg.From, txData)
	}

	if err != nil {
		return nil, err
	}

	return &sig, nil
}

func (i *Interceptor) EthSignTransaction() jsonrpc.Handler {
	h, _ := jsonrpc.MakeHandler(i.ethSignTransaction)
	return h
}
