package interceptor

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/ethereum"
	"github.com/consensys/quorum-key-manager/pkg/jsonrpc"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (i *Interceptor) eeaSendTransaction(ctx context.Context, msg *ethereum.SendEEATxMsg) (*ethcommon.Hash, error) {
	i.logger.Debug("sending EEA transaction")

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	// Get store for from
	store, err := i.stores.EthereumByAddr(ctx, msg.From, userInfo)
	if err != nil {
		return nil, err
	}

	sess := proxynode.SessionFromContext(ctx)

	if msg.Nonce == nil {
		var n uint64
		if msg.PrivacyGroupID != nil {
			n, err = sess.EthCaller().Priv().GetTransactionCount(ctx, msg.From, *msg.PrivacyGroupID)
		} else {
			if msg.PrivateFor == nil {
				errMessage := "missing privateFor"
				i.logger.Error(errMessage)
				return nil, errors.InvalidFormatError(errMessage)
			}

			var privateFrom string
			if msg.PrivateFrom != nil {
				privateFrom = *msg.PrivateFrom
			}
			n, err = sess.EthCaller().Priv().GetEeaTransactionCount(ctx, msg.From, privateFrom, *msg.PrivateFor)
		}
		if err != nil {
			i.logger.WithError(err).Error("failed to fetch transaction count (EEA transaction)")
			return nil, errors.BlockchainNodeError(err.Error())
		}

		msg.Nonce = &n
	}

	if msg.GasPrice == nil {
		gasPrice, err2 := sess.EthCaller().Eth().GasPrice(ctx)
		if err2 != nil {
			i.logger.WithError(err2).Error("failed to fetch gas price (EEA transaction)")
			return nil, errors.BlockchainNodeError(err2.Error())
		}

		msg.GasPrice = gasPrice
	}

	if msg.Gas == nil {
		// We update the data to an arbitrary hash
		// to avoid errors raised on eth_estimateGas on Besu 1.5.4 & 1.5.5
		callMsg := &ethereum.CallMsg{
			From:     &msg.From,
			To:       msg.To,
			GasPrice: msg.GasPrice,
			Data:     common.ToPtr(hexutil.MustDecode("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")).(*[]byte),
		}

		gas, err2 := sess.EthCaller().Eth().EstimateGas(ctx, callMsg)
		if err2 != nil {
			i.logger.WithError(err2).Error("failed to estimate gas (EEA transaction)")
			return nil, errors.BlockchainNodeError(err2.Error())
		}

		msg.Gas = &gas
	}

	if msg.PrivateType == nil {
		msg.PrivateType = common.ToPtr(ethereum.PrivateTypeRestricted).(*ethereum.PrivateType)
	}

	// Get ChainID from Node
	chainID, err := sess.EthCaller().Eth().ChainID(ctx)
	if err != nil {
		i.logger.WithError(err).Error("failed to fetch chainID (EEA transaction)")
		return nil, errors.BlockchainNodeError(err.Error())
	}

	// Sign
	sig, err := store.SignEEA(ctx, msg.From, chainID, msg.TxData(), &msg.PrivateArgs)
	if err != nil {
		return nil, err
	}

	// Submit transaction to downstream node
	hash, err := sess.EthCaller().EEA().SendRawTransaction(ctx, sig)
	if err != nil {
		i.logger.WithError(err).Error("failed to send raw EEA transaction")
		return nil, errors.BlockchainNodeError(err.Error())
	}

	i.logger.Info("EEA transaction sent successfully", "tx_hash", hash)
	return &hash, nil
}

func (i *Interceptor) EEASendTransaction() jsonrpc.Handler {
	h, _ := jsonrpc.MakeHandler(i.eeaSendTransaction)
	return h
}
