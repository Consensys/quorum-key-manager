package ethereum

import (
	"context"
	"math/big"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func init() {
	err := jsonrpc.ProvideCaller(
		ethSrv,
		eeaSrv,
		privSrv,
	)
	if err != nil {
		panic(err)
	}
}

//go:generate mockgen -source=caller.go -destination=mock/caller.go -package=mock

type Caller interface {
	Eth() EthCaller
	EEA() EEACaller
	Priv() PrivCaller
}

type EthCaller interface {
	ChainID(context.Context) (*big.Int, error)
	GasPrice(context.Context) (*big.Int, error)
	GetTransactionCount(context.Context, ethcommon.Address, BlockNumber) (uint64, error)
	EstimateGas(context.Context, *CallMsg) (uint64, error)
	SendRawTransaction(context.Context, []byte) (ethcommon.Hash, error)
	SendRawPrivateTransaction(context.Context, []byte, *PrivateArgs) (ethcommon.Hash, error)
}

type EEACaller interface {
	SendRawTransaction(context.Context, []byte) (ethcommon.Hash, error)
}

type PrivCaller interface {
	DistributeRawTransaction(context.Context, []byte) ([]byte, error)
	GetTransactionCount(ctx context.Context, addr ethcommon.Address, privacyGroupID string) (uint64, error)
	GetEEATransactionCount(ctx context.Context, addr ethcommon.Address, privateFrom string, privateFor []string) (uint64, error)
}

// Caller implement methods to interface with the JSON-RPC API of an ethereum caller
// It is some kind of Web3 interface
type caller struct {
	eth  *ethCaller
	eea  *eeaCaller
	priv *privCaller
}

// NewCaller creates a caller from a jsonrpc.Client
func NewCaller(client jsonrpc.Client) Caller {
	return &caller{
		eth:  &ethCaller{client},
		eea:  &eeaCaller{client},
		priv: &privCaller{client},
	}
}

// Eth return eth namespace caller
func (c *caller) Eth() EthCaller { // nolint
	return c.eth
}

// EEA return eea namespace caller
func (c *caller) EEA() EEACaller { // nolint
	return c.eea
}

// Priv return priv namespace caller
func (c *caller) Priv() PrivCaller { // nolint
	return c.priv
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
	return ethSrv.SendRawTransaction(c.client)(ctx, hexutil.Bytes(raw))
}

func (c *ethCaller) SendRawPrivateTransaction(ctx context.Context, raw []byte, privArgs *PrivateArgs) (ethcommon.Hash, error) {
	return ethSrv.SendRawPrivateTransaction(c.client)(ctx, hexutil.Bytes(raw), privArgs)
}

type eeaCaller struct {
	client jsonrpc.Client
}

func (c *eeaCaller) SendRawTransaction(ctx context.Context, raw []byte) (ethcommon.Hash, error) {
	return eeaSrv.SendRawTransaction(c.client)(ctx, hexutil.Bytes(raw))
}

type privCaller struct {
	client jsonrpc.Client
}

func (c *privCaller) DistributeRawTransaction(ctx context.Context, raw []byte) ([]byte, error) {
	enclaveKey, err := privSrv.DistributeRawTransaction(c.client)(ctx, hexutil.Bytes(raw))
	if err != nil {
		return nil, err
	}

	return []byte(*enclaveKey), nil
}

func (c *privCaller) GetTransactionCount(ctx context.Context, addr ethcommon.Address, privacyGroupID string) (uint64, error) {
	n, err := privSrv.GetTransactionCount(c.client)(ctx, addr, privacyGroupID)
	if err != nil {
		return 0, err
	}

	return uint64(*n), nil
}

func (c *privCaller) GetEEATransactionCount(ctx context.Context, addr ethcommon.Address, privateFrom string, privateFor []string) (uint64, error) {
	n, err := privSrv.GetEEATransactionCount(c.client)(ctx, addr, privateFrom, privateFor)
	if err != nil {
		return 0, err
	}

	return uint64(*n), nil
}

var (
	ethSrv  = new(ethService)
	eeaSrv  = new(eeaService)
	privSrv = new(privService)
)

type ethService struct {
	ChainID                   func(jsonrpc.Client) func(context.Context) (*hexutil.Big, error)                                    `method:"eth_chainId"`
	GasPrice                  func(jsonrpc.Client) func(context.Context) (*hexutil.Big, error)                                    `namespace:"eth"`
	GetTransactionCount       func(jsonrpc.Client) func(context.Context, ethcommon.Address, BlockNumber) (*hexutil.Uint64, error) `namespace:"eth"`
	EstimateGas               func(jsonrpc.Client) func(context.Context, *CallMsg) (*hexutil.Uint64, error)                       `namespace:"eth"`
	SendRawTransaction        func(jsonrpc.Client) func(context.Context, hexutil.Bytes) (ethcommon.Hash, error)                   `namespace:"eth"`
	SendRawPrivateTransaction func(jsonrpc.Client) func(context.Context, hexutil.Bytes, *PrivateArgs) (ethcommon.Hash, error)     `namespace:"eth"`
}

type eeaService struct {
	SendRawTransaction func(jsonrpc.Client) func(context.Context, hexutil.Bytes) (ethcommon.Hash, error) `namespace:"eea"`
}

type privService struct {
	DistributeRawTransaction func(jsonrpc.Client) func(context.Context, hexutil.Bytes) (*hexutil.Bytes, error)                        `namespace:"priv"`
	GetEEATransactionCount   func(jsonrpc.Client) func(context.Context, ethcommon.Address, string, []string) (*hexutil.Uint64, error) `namespace:"priv"`
	GetTransactionCount      func(jsonrpc.Client) func(context.Context, ethcommon.Address, string) (*hexutil.Uint64, error)           `namespace:"priv"`
}
