package interceptor

import (
	"context"

	"github.com/consensysquorum/quorum-key-manager/pkg/common"
	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	"github.com/consensysquorum/quorum-key-manager/pkg/ethereum"
	"github.com/consensysquorum/quorum-key-manager/pkg/jsonrpc"
	proxynode "github.com/consensysquorum/quorum-key-manager/src/nodes/node/proxy"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	privateTxTypeRestricted = "restricted"
)

func (i *Interceptor) eeaSendTransaction(ctx context.Context, msg *ethereum.SendEEATxMsg) (*ethcommon.Hash, error) {
	i.logger.Debug("sending EEA transaction")

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
			return nil, err
		}

		msg.Nonce = &n
	}

	if msg.GasPrice == nil {
		gasPrice, err2 := sess.EthCaller().Eth().GasPrice(ctx)
		if err2 != nil {
			errMessage := "failed to fetch gas price (EEA transaction)"
			i.logger.WithError(err2).Error(errMessage)
			return nil, errors.BlockchainNodeError(errMessage)
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
			errMessage := "failed to estimate gas (EEA transaction)"
			i.logger.WithError(err2).Error(errMessage)
			return nil, errors.BlockchainNodeError(errMessage)
		}

		msg.Gas = &gas
	}

	if msg.PrivateType == nil {
		msg.PrivateType = common.ToPtr(privateTxTypeRestricted).(*string)
	}

	// Get ChainID from Node
	chainID, err := sess.EthCaller().Eth().ChainID(ctx)
	if err != nil {
		errMessage := "failed to fetch chainID (EEA transaction)"
		i.logger.WithError(err).Error(errMessage)
		return nil, errors.BlockchainNodeError(errMessage)
	}

	// Sign
	sig, err := store.SignEEA(ctx, msg.From.Hex(), chainID, msg.TxData(), &msg.PrivateArgs)
	if err != nil {
		return nil, err
	}

	// Submit transaction to downstream node
	hash, err := sess.EthCaller().EEA().SendRawTransaction(ctx, sig)
	if err != nil {
		errMessage := "failed to send raw EEA transaction"
		i.logger.WithError(err).Error(errMessage)
		return nil, errors.BlockchainNodeError(errMessage)
	}

	i.logger.Info("EEA transaction sent successfully", "tx_hash", hash)
	return &hash, nil
}

func (i *Interceptor) EEASendTransaction() jsonrpc.Handler {
	h, _ := jsonrpc.MakeHandler(i.eeaSendTransaction)
	return h
}
