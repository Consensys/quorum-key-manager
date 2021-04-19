package ethereum

import (
	"context"

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

// Client implement methods to interface with the JSON-RPC API of an ethereum client
// It is some kind of Web3 interface
type Client struct {
	eth  *ethClient
	eea  *eeaClient
	priv *privClient
}

// NewClient creates a client from a jsonrpc.Caller
func NewClient(cllr jsonrpc.Caller) *Client {
	return &Client{
		eth:  &ethClient{cllr},
		eea:  &eeaClient{cllr},
		priv: &privClient{cllr},
	}
}

// Eth return eth namespace client
func (c *Client) Eth() *ethClient { // nolint
	return c.eth
}

// EEA return eea namespace client
func (c *Client) EEA() *eeaClient { // nolint
	return c.eea
}

// Priv return priv namespace client
func (c *Client) Priv() *privClient { // nolint
	return c.priv
}

type ethClient struct {
	cllr jsonrpc.Caller
}

func (c *ethClient) ChainID(ctx context.Context) (*hexutil.Big, error) {
	return ethSrv.ChainID(c.cllr)(ctx)
}

func (c *ethClient) GasPrice(ctx context.Context) (*hexutil.Big, error) {
	return ethSrv.GasPrice(c.cllr)(ctx)
}

func (c *ethClient) GetTransactionCount(ctx context.Context, addr ethcommon.Address, blockNumber BlockNumber) (*hexutil.Uint64, error) {
	return ethSrv.GetTransactionCount(c.cllr)(ctx, addr, blockNumber)
}

func (c *ethClient) EstimateGas(ctx context.Context, msg *CallMsg) (*hexutil.Uint64, error) {
	return ethSrv.EstimateGas(c.cllr)(ctx, msg)
}

func (c *ethClient) SendRawTransaction(ctx context.Context, raw hexutil.Bytes) (ethcommon.Hash, error) {
	return ethSrv.SendRawTransaction(c.cllr)(ctx, raw)
}

func (c *ethClient) SendRawPrivateTransaction(ctx context.Context, raw hexutil.Bytes, privArgs *PrivateArgs) (ethcommon.Hash, error) {
	return ethSrv.SendRawPrivateTransaction(c.cllr)(ctx, raw, privArgs)
}

type eeaClient struct {
	cllr jsonrpc.Caller
}

func (c *eeaClient) SendRawTransaction(ctx context.Context, raw hexutil.Bytes) (ethcommon.Hash, error) {
	return eeaSrv.SendRawTransaction(c.cllr)(ctx, raw)
}

type privClient struct {
	cllr jsonrpc.Caller
}

func (c *privClient) DistributeRawTransaction(ctx context.Context, raw hexutil.Bytes) (*hexutil.Bytes, error) {
	return privSrv.DistributeRawTransaction(c.cllr)(ctx, raw)
}

func (c *privClient) GetTransactionCount(ctx context.Context, addr ethcommon.Address, privacyGroupID string) (*hexutil.Uint64, error) {
	return privSrv.GetTransactionCount(c.cllr)(ctx, addr, privacyGroupID)
}

func (c *privClient) GetEEATransactionCount(ctx context.Context, addr ethcommon.Address, privateFrom string, privateFor []string) (*hexutil.Uint64, error) {
	return privSrv.GetEEATransactionCount(c.cllr)(ctx, addr, privateFrom, privateFor)
}

var (
	ethSrv  = new(ethService)
	eeaSrv  = new(eeaService)
	privSrv = new(privService)
)

type ethService struct {
	ChainID                   func(jsonrpc.Caller) func(context.Context) (*hexutil.Big, error)                                    `method:"eth_chainId"`
	GasPrice                  func(jsonrpc.Caller) func(context.Context) (*hexutil.Big, error)                                    `namespace:"eth"`
	GetTransactionCount       func(jsonrpc.Caller) func(context.Context, ethcommon.Address, BlockNumber) (*hexutil.Uint64, error) `namespace:"eth"`
	EstimateGas               func(jsonrpc.Caller) func(context.Context, *CallMsg) (*hexutil.Uint64, error)                       `namespace:"eth"`
	SendRawTransaction        func(jsonrpc.Caller) func(context.Context, hexutil.Bytes) (ethcommon.Hash, error)                   `namespace:"eth"`
	SendRawPrivateTransaction func(jsonrpc.Caller) func(context.Context, hexutil.Bytes, *PrivateArgs) (ethcommon.Hash, error)     `namespace:"eth"`
}

type eeaService struct {
	SendRawTransaction func(jsonrpc.Caller) func(context.Context, hexutil.Bytes) (ethcommon.Hash, error) `namespace:"eea"`
}

type privService struct {
	DistributeRawTransaction func(jsonrpc.Caller) func(context.Context, hexutil.Bytes) (*hexutil.Bytes, error)                        `namespace:"priv"`
	GetEEATransactionCount   func(jsonrpc.Caller) func(context.Context, ethcommon.Address, string, []string) (*hexutil.Uint64, error) `namespace:"priv"`
	GetTransactionCount      func(jsonrpc.Caller) func(context.Context, ethcommon.Address, string) (*hexutil.Uint64, error)           `namespace:"priv"`
}
