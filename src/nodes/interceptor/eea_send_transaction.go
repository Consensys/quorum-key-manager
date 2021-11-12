package interceptor

import (
	"context"
	"fmt"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/ethereum"
	"github.com/consensys/quorum-key-manager/pkg/jsonrpc"
	"github.com/consensys/quorum-key-manager/src/aliases/entities"
	"github.com/consensys/quorum-key-manager/src/auth/api/middlewares"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (i *Interceptor) eeaSendTransaction(ctx context.Context, msg *ethereum.SendEEATxMsg) (*ethcommon.Hash, error) {
	i.logger.Debug("sending EEA transaction")

	// Get store for from
	store, err := i.stores.EthereumByAddr(ctx, msg.From, middlewares.UserInfoFromContext(ctx))
	if err != nil {
		return nil, err
	}

	sess := proxynode.SessionFromContext(ctx)

	if msg.PrivateFor != nil {
		// extract aliases from PrivateFor
		*msg.PrivateFor, err = i.aliases.ReplaceAliases(ctx, *msg.PrivateFor)
		if err != nil {
			i.logger.WithError(err).Error("failed to replace aliases in privateFor")
			return nil, err
		}
	}

	if msg.PrivateFrom != nil {

		*msg.PrivateFrom, err = i.aliases.ReplaceSimpleAlias(ctx, *msg.PrivateFrom)
		if err != nil {
			i.logger.WithError(err).Error("failed to replace alias")
			return nil, err
		}
	}

	if msg.PrivacyGroupID != nil {
		reg, key, isAlias := i.aliases.ParseAlias(*msg.PrivacyGroupID)
		if isAlias {
			var alias *entities.Alias
			alias, err = i.aliases.GetAlias(ctx, reg, key)
			if err != nil {
				i.logger.WithError(err).Error("failed to get alias for privacyGroupID")
				return nil, err
			}

			switch alias.Value.Kind {
			case entities.KindString:
				*msg.PrivacyGroupID, err = alias.Value.String()
				if err != nil {
					i.logger.WithError(err).Error("wrong alias value, should be a string")
					return nil, err
				}
			case entities.KindArray:
				if msg.PrivateFor == nil {
					slice := []string{}
					msg.PrivateFor = &slice
				}

				var aliasArray []string
				aliasArray, err = alias.Value.Array()
				if err != nil {
					i.logger.WithError(err).Error("wrong alias value, should be a string")
					return nil, err
				}

				*msg.PrivateFor = append(*msg.PrivateFor, aliasArray...)
				msg.PrivacyGroupID = nil
			default:
				msg := "wrong alias type"
				err = fmt.Errorf(msg)
				i.logger.WithError(err).Error(msg)
				return nil, err
			}
		}
	}

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
