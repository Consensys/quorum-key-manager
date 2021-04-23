package interceptor

import (
	"context"
	"math/big"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/node"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (i *Interceptor) eeaSendTransaction(ctx context.Context, msg *ethereum.SendTxMsg) (*ethcommon.Hash, error) {
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
	sig, err = store.SignEEA(ctx, (*big.Int)(chainID), msg.From, txData, &msg.PrivateArgs)
	if err != nil {
		return nil, err
	}

	// Submit transaction to downstream node
	hash, err := sess.EthClient().EEA().SendRawTransaction(ctx, sig)
	if err != nil {
		return nil, err
	}

	return &hash, nil
}

func (i *Interceptor) EEASendTransaction() jsonrpc.Handler {
	h, _ := jsonrpc.MakeHandler(i.eeaSendTransaction)
	return h
}
