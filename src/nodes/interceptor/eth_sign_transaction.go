package interceptor

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator"

	"github.com/consensys/quorum-key-manager/pkg/ethereum"
	"github.com/consensys/quorum-key-manager/pkg/jsonrpc"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (i *Interceptor) ethSignTransaction(ctx context.Context, msg *ethereum.SendTxMsg) (*hexutil.Bytes, error) {
	i.logger.Debug("signing ETH transaction")

	if msg.Gas == nil {
		errMessage := "gas not specified"
		i.logger.Error(errMessage)
		return nil, jsonrpc.InvalidParamsError(errors.InvalidParameterError(errMessage))
	}

	if msg.GasPrice == nil {
		errMessage := "gasPrice not specified"
		i.logger.Error(errMessage)
		return nil, jsonrpc.InvalidParamsError(errors.InvalidParameterError(errMessage))
	}

	if msg.Nonce == nil {
		errMessage := "nonce not specified"
		i.logger.Error(errMessage)
		return nil, jsonrpc.InvalidParamsError(errors.InvalidParameterError(errMessage))
	}

	if msg.Data == nil {
		msg.Data = &[]byte{}
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	// Get store for from
	store, err := i.stores.GetEth1StoreByAddr(ctx, msg.From, userInfo)
	if err != nil {
		return nil, err
	}

	// Get ChainID from Node
	sess := proxynode.SessionFromContext(ctx)
	chainID, err := sess.EthCaller().Eth().ChainID(ctx)
	if err != nil {
		i.logger.WithError(err).Error("failed to fetch chainID")
		return nil, errors.BlockchainNodeError(err.Error())
	}

	// Sign
	var sig []byte
	if msg.IsPrivate() {
		sig, err = store.SignPrivate(ctx, msg.From.Hex(), msg.TxDataQuorum())
	} else {
		sig, err = store.SignTransaction(ctx, msg.From.Hex(), chainID, msg.TxData())
	}
	if err != nil {
		return nil, err
	}

	i.logger.Info("ETH transaction signed successfully")
	return (*hexutil.Bytes)(&sig), nil
}

func (i *Interceptor) EthSignTransaction() jsonrpc.Handler {
	h, _ := jsonrpc.MakeHandler(i.ethSignTransaction)
	return h
}
