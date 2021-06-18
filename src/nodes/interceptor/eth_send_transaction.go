package interceptor

import (
	"context"

	"github.com/consensysquorum/quorum-key-manager/pkg/errors"

	"github.com/consensysquorum/quorum-key-manager/pkg/ethereum"
	"github.com/consensysquorum/quorum-key-manager/pkg/jsonrpc"
	proxynode "github.com/consensysquorum/quorum-key-manager/src/nodes/node/proxy"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func (i *Interceptor) ethSendTransaction(ctx context.Context, msg *ethereum.SendTxMsg) (*ethcommon.Hash, error) {
	i.logger.Debug("sending ETH transaction")

	// Get ChainID from Node
	sess := proxynode.SessionFromContext(ctx)

	if msg.GasPrice == nil {
		gasPrice, err := sess.EthCaller().Eth().GasPrice(ctx)
		if err != nil {
			i.logger.WithError(err).Error("failed to fetch gas price")
			return nil, errors.BlockchainNodeError(err.Error())
		}

		msg.GasPrice = gasPrice
	}

	if msg.Gas == nil {
		callMsg := &ethereum.CallMsg{
			From:     &msg.From,
			To:       msg.To,
			GasPrice: msg.GasPrice,
			Value:    msg.Value,
			Data:     msg.Data,
		}
		gas, err := sess.EthCaller().Eth().EstimateGas(ctx, callMsg)
		if err != nil {
			i.logger.WithError(err).Error("failed to estimate gas")
			return nil, errors.BlockchainNodeError(err.Error())
		}

		msg.Gas = &gas
	}

	if msg.Nonce == nil {
		n, err := sess.EthCaller().Eth().GetTransactionCount(ctx, msg.From, ethereum.PendingBlockNumber)
		if err != nil {
			i.logger.WithError(err).Error("failed to fetch nonce", "from_account", msg.From)
			return nil, errors.BlockchainNodeError(err.Error())
		}

		msg.Nonce = &n
	}

	if msg.IsPrivate() {
		if msg.Data == nil {
			msg.Data = new([]byte)
		}

		var privateFrom string
		if msg.PrivateFrom != nil {
			privateFrom = *msg.PrivateFrom
		}

		// Store payload on Tessera
		key, err := sess.ClientPrivTxManager().StoreRaw(ctx, *msg.Data, privateFrom)
		if err != nil {
			i.logger.WithError(err).Error("failed to store raw payload on Tessera", "private_from", privateFrom)
			return nil, errors.BlockchainNodeError(err.Error())
		}

		// Switch message data
		*msg.Data = key
	}

	raw, err := i.ethSignTransaction(ctx, msg)
	if err != nil {
		return nil, err
	}

	var hash ethcommon.Hash
	if msg.IsPrivate() {
		hash, err = sess.EthCaller().Eth().SendRawPrivateTransaction(ctx, *raw, &msg.PrivateArgs)
	} else {
		hash, err = sess.EthCaller().Eth().SendRawTransaction(ctx, *raw)
	}

	if err != nil {
		i.logger.WithError(err).Error("failed to store raw transaction")
		return nil, errors.BlockchainNodeError(err.Error())
	}

	i.logger.Info("ETH transaction sent successfully", "tx_hash", hash)
	return &hash, nil
}

func (i *Interceptor) EthSendTransaction() jsonrpc.Handler {
	h, _ := jsonrpc.MakeHandler(i.ethSendTransaction)
	return h
}
