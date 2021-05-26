package interceptor

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	proxynode "github.com/ConsenSysQuorum/quorum-key-manager/src/services/nodes/node/proxy"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func (i *Interceptor) eeaSendTransaction(ctx context.Context, msg *ethereum.SendEEATxMsg) (*ethcommon.Hash, error) {
	// Get store for from
	store, err := i.stores.GetEth1StoreByAddr(ctx, msg.From)
	if err != nil {
		return nil, err
	}

	sess := proxynode.SessionFromContext(ctx)

	if msg.Nonce == nil {
		var n uint64
		if msg.PrivacyGroupID != nil {
			n, err = sess.EthCaller().Priv().GetTransactionCount(ctx, msg.From, *msg.PrivacyGroupID)
		} else {
			var privateFrom string
			if msg.PrivateFrom != nil {
				privateFrom = *msg.PrivateFrom
			}
			n, err = sess.EthCaller().Priv().GetEEATransactionCount(ctx, msg.From, privateFrom, *msg.PrivateFor)
		}

		if err != nil {
			return nil, err
		}

		msg.Nonce = &n
	}

	// Get ChainID from Node
	chainID, err := sess.EthCaller().Eth().ChainID(ctx)
	if err != nil {
		return nil, err
	}

	// Sign
	sig, err := store.SignEEA(ctx, msg.From.Hex(), chainID, msg.TxData(), &msg.PrivateArgs)
	if err != nil {
		return nil, err
	}

	// Submit transaction to downstream node
	hash, err := sess.EthCaller().EEA().SendRawTransaction(ctx, sig)
	if err != nil {
		return nil, err
	}

	return &hash, nil
}

func (i *Interceptor) EEASendTransaction() jsonrpc.Handler {
	h, _ := jsonrpc.MakeHandler(i.eeaSendTransaction)
	return h
}
