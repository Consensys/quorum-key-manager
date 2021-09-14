package interceptor

import (
	"context"
	"math/big"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/ethereum"
	"github.com/consensys/quorum-key-manager/pkg/jsonrpc"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func (i *Interceptor) ethSendTransaction(ctx context.Context, msg *ethereum.SendTxMsg) (*ethcommon.Hash, error) {
	switch {
	case msg.IsPrivate():
		return i.sendPrivateTx(ctx, msg)
	case msg.GasPrice != nil:
		return i.sendLegacyTx(ctx, msg)
	default:
		return i.sendTx(ctx, msg)
	}
}

func (i *Interceptor) sendPrivateTx(ctx context.Context, msg *ethereum.SendTxMsg) (*ethcommon.Hash, error) {
	i.logger.Debug("sending Quorum private transaction")

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
			i.logger.WithError(err).With("gas_price", msg.GasPrice).Error("failed to estimate gas for quorum private transaction")
			return nil, errors.BlockchainNodeError(err.Error())
		}

		msg.Gas = &gas
	}

	err := i.fillNonce(ctx, sess, msg)
	if err != nil {
		return nil, err
	}

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

	raw, err := i.ethSignTransaction(ctx, msg)
	if err != nil {
		return nil, err
	}

	hash, err := sess.EthCaller().Eth().SendRawPrivateTransaction(ctx, *raw, &msg.PrivateArgs)
	if err != nil {
		i.logger.WithError(err).Error("failed to send raw quorum private transaction")
		return nil, errors.BlockchainNodeError(err.Error())
	}

	i.logger.Info("quorum private transaction sent successfully", "tx_hash", hash)
	return &hash, nil
}

func (i *Interceptor) sendLegacyTx(ctx context.Context, msg *ethereum.SendTxMsg) (*ethcommon.Hash, error) {
	i.logger.Debug("sending ETH legacy transaction")

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
			i.logger.WithError(err).With("gas_price", msg.GasPrice).Error("failed to estimate gas for legacy transaction")
			return nil, errors.BlockchainNodeError(err.Error())
		}

		msg.Gas = &gas
	}

	err := i.fillNonce(ctx, sess, msg)
	if err != nil {
		return nil, err
	}

	raw, err := i.ethSignTransaction(ctx, msg)
	if err != nil {
		return nil, err
	}

	hash, err := sess.EthCaller().Eth().SendRawTransaction(ctx, *raw)
	if err != nil {
		i.logger.WithError(err).Error("failed to send raw legacy transaction")
		return nil, errors.BlockchainNodeError(err.Error())
	}

	i.logger.Info("legacy transaction sent successfully", "tx_hash", hash)
	return &hash, nil
}

func (i *Interceptor) sendTx(ctx context.Context, msg *ethereum.SendTxMsg) (*ethcommon.Hash, error) {
	i.logger.Debug("sending ETH transaction")

	sess := proxynode.SessionFromContext(ctx)

	baseFee, err := sess.EthCaller().Eth().BaseFeePerGas(ctx, ethereum.LatestBlockNumber)
	if err != nil {
		i.logger.WithError(err).Error("failed to retrieve base fee from latest block")
		return nil, errors.BlockchainNodeError(err.Error())
	}

	if baseFee == nil {
		i.logger.Warn("cannot send a dynamic fee transaction to a pre-London node, reverting to legacy tx")
		return i.sendLegacyTx(ctx, msg)
	}

	if msg.GasFeeCap == nil {
		var maxPriorityFeePerGas = big.NewInt(0)
		if msg.GasTipCap != nil {
			maxPriorityFeePerGas = msg.GasTipCap
		}
		msg.GasFeeCap = new(big.Int).Add(baseFee, maxPriorityFeePerGas)
		i.logger.
			With("max_fee_per_gas", msg.GasFeeCap, "base_fee", baseFee, "max_priority_fee_per_gas", maxPriorityFeePerGas).
			Debug("'maxFeePerGas' set with previous block 'baseFeePerGas' + miner tip")
	}

	if msg.Gas == nil {
		callMsg := &ethereum.CallMsg{
			From:       &msg.From,
			To:         msg.To,
			Value:      msg.Value,
			Data:       msg.Data,
			GasTipCap:  msg.GasTipCap,
			GasFeeCap:  msg.GasFeeCap,
			AccessList: msg.AccessList,
		}
		gas, err2 := sess.EthCaller().Eth().EstimateGas(ctx, callMsg)
		if err2 != nil {
			i.logger.WithError(err2).With("gas_price", msg.GasPrice).Error("failed to estimate gas for dynamic fee transaction")
			return nil, errors.BlockchainNodeError(err2.Error())
		}

		msg.Gas = &gas
	}

	err = i.fillNonce(ctx, sess, msg)
	if err != nil {
		return nil, err
	}

	raw, err := i.ethSignTransaction(ctx, msg)
	if err != nil {
		return nil, err
	}

	hash, err := sess.EthCaller().Eth().SendRawTransaction(ctx, *raw)
	if err != nil {
		i.logger.WithError(err).Error("failed to send raw transaction")
		return nil, errors.BlockchainNodeError(err.Error())
	}

	i.logger.Info("ETH transaction sent successfully", "tx_hash", hash)
	return &hash, nil
}

func (i *Interceptor) fillNonce(ctx context.Context, sess proxynode.Session, msg *ethereum.SendTxMsg) error {
	if msg.Nonce == nil {
		n, err := sess.EthCaller().Eth().GetTransactionCount(ctx, msg.From, ethereum.PendingBlockNumber)
		if err != nil {
			i.logger.WithError(err).Error("failed to fetch nonce", "from_account", msg.From)
			return errors.BlockchainNodeError(err.Error())
		}

		msg.Nonce = &n
	}

	return nil
}

func (i *Interceptor) EthSendTransaction() jsonrpc.Handler {
	h, _ := jsonrpc.MakeHandler(i.ethSendTransaction)
	return h
}
