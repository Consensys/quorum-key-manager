package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/consensys/quorum-key-manager/pkg/jsonrpc"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func init() {
	err := jsonrpc.ProvideCaller(
		ethSrv,
	)
	if err != nil {
		panic(err)
	}
}

var ethSrv = new(ethService)

// ethService is a jsonrpc.Caller which methods are meant to be automatically populated using jsonrpc.ProvideCaller
type ethService struct {
	ChainID                   func(jsonrpc.Client) func(context.Context) (*hexutil.Big, error)                                    `method:"eth_chainId"`
	GasPrice                  func(jsonrpc.Client) func(context.Context) (*hexutil.Big, error)                                    `namespace:"eth"`
	GetTransactionCount       func(jsonrpc.Client) func(context.Context, ethcommon.Address, BlockNumber) (*hexutil.Uint64, error) `namespace:"eth"`
	EstimateGas               func(jsonrpc.Client) func(context.Context, *CallMsg) (*hexutil.Uint64, error)                       `namespace:"eth"`
	SendRawTransaction        func(jsonrpc.Client) func(context.Context, hexutil.Bytes) (ethcommon.Hash, error)                   `namespace:"eth"`
	SendRawPrivateTransaction func(jsonrpc.Client) func(context.Context, hexutil.Bytes, *PrivateArgs) (ethcommon.Hash, error)     `namespace:"eth"`
	GetBlockByNumber          func(jsonrpc.Client) func(context.Context, BlockNumber, bool) (*types.Header, error)                `method:"eth_getBlockByNumber"`
}

//go:generate mockgen -source=caller_eth.go -destination=mock/caller_eth.go -package=mock

// EthCaller is a JSON-RPC client to an Ethereum client using eth namespace
type EthCaller interface {
	ChainID(context.Context) (*big.Int, error)
	GasPrice(context.Context) (*big.Int, error)
	BaseFeePerGas(context.Context, BlockNumber) (*big.Int, error)
	GetTransactionCount(context.Context, ethcommon.Address, BlockNumber) (uint64, error)
	EstimateGas(context.Context, *CallMsg) (uint64, error)
	SendRawTransaction(context.Context, []byte) (ethcommon.Hash, error)
	SendRawPrivateTransaction(context.Context, []byte, *PrivateArgs) (ethcommon.Hash, error)
}

type ethCaller struct {
	client jsonrpc.Client
}

func (c *ethCaller) ChainID(ctx context.Context) (*big.Int, error) {
	chainID, err := ethSrv.ChainID(c.client)(ctx)
	if err != nil {
		return nil, err
	}

	return (*big.Int)(chainID), nil
}

func (c *ethCaller) GasPrice(ctx context.Context) (*big.Int, error) {
	p, err := ethSrv.GasPrice(c.client)(ctx)
	if err != nil {
		return nil, err
	}

	return (*big.Int)(p), nil
}

func (c *ethCaller) GetTransactionCount(ctx context.Context, addr ethcommon.Address, blockNumber BlockNumber) (uint64, error) {
	n, err := ethSrv.GetTransactionCount(c.client)(ctx, addr, blockNumber)
	if err != nil {
		return 0, err
	}

	return uint64(*n), nil
}

func (c *ethCaller) EstimateGas(ctx context.Context, msg *CallMsg) (uint64, error) {
	gas, err := ethSrv.EstimateGas(c.client)(ctx, msg)
	if err != nil {
		return 0, err
	}

	return uint64(*gas), nil
}

func (c *ethCaller) SendRawTransaction(ctx context.Context, raw []byte) (ethcommon.Hash, error) {
	return ethSrv.SendRawTransaction(c.client)(ctx, raw)
}

func (c *ethCaller) SendRawPrivateTransaction(ctx context.Context, raw []byte, privArgs *PrivateArgs) (ethcommon.Hash, error) {
	return ethSrv.SendRawPrivateTransaction(c.client)(ctx, raw, privArgs)
}

func (c *ethCaller) BaseFeePerGas(ctx context.Context, blockNumber BlockNumber) (*big.Int, error) {
	header, err := ethSrv.GetBlockByNumber(c.client)(ctx, blockNumber, false)
	if err != nil {
		return nil, err
	}

	return header.BaseFee, nil
}
