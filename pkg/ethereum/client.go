package ethereum

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func init() {
	err := jsonrpc.Provide(
		ethSrv,
		eeaSrv,
		privSrv,
	)
	if err != nil {
		panic(err)
	}
}

type Client struct {
	eth  *ethClient
	eea  *eeaClient
	priv *privClient
}

func NewClient(cllr jsonrpc.Caller) *Client {
	return &Client{
		eth:  &ethClient{cllr},
		eea:  &eeaClient{cllr},
		priv: &privClient{cllr},
	}
}

func (c *Client) Eth() *ethClient { // nolint
	return c.eth
}

func (c *Client) EEA() *eeaClient { // nolint
	return c.eea
}

func (c *Client) Priv() *privClient { // nolint
	return c.priv
}

type ethClient struct {
	cllr jsonrpc.Caller
}

func (c *ethClient) GasPrice(ctx context.Context) (*hexutil.Big, error) {
	return ethSrv.GasPrice(c.cllr)(ctx)
}

func (c *ethClient) GetTransactionCount(ctx context.Context, addr ethcommon.Address, blockNumber BlockNumber) (*hexutil.Uint64, error) {
	return ethSrv.GetTransactionCount(c.cllr)(ctx, addr, blockNumber)
}

func (c *ethClient) SendRawTransaction(ctx context.Context, raw hexutil.Bytes) (ethcommon.Hash, error) {
	return ethSrv.SendRawTransaction(c.cllr)(ctx, raw)
}

func (c *ethClient) EstimateGas(ctx context.Context, msg *CallMsg) (*hexutil.Uint64, error) {
	return ethSrv.EstimateGas(c.cllr)(ctx, msg)
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
	GasPrice            func(jsonrpc.Caller) func(context.Context) (*hexutil.Big, error)                                    `namespace:"eth"`
	GetTransactionCount func(jsonrpc.Caller) func(context.Context, ethcommon.Address, BlockNumber) (*hexutil.Uint64, error) `namespace:"eth"`
	SendRawTransaction  func(jsonrpc.Caller) func(context.Context, hexutil.Bytes) (ethcommon.Hash, error)                   `namespace:"eth"`
	EstimateGas         func(jsonrpc.Caller) func(context.Context, *CallMsg) (*hexutil.Uint64, error)                       `namespace:"eth"`
}

type eeaService struct {
	SendRawTransaction func(jsonrpc.Caller) func(context.Context, hexutil.Bytes) (ethcommon.Hash, error) `namespace:"eea"`
}

type privService struct {
	DistributeRawTransaction func(jsonrpc.Caller) func(context.Context, hexutil.Bytes) (*hexutil.Bytes, error)                        `namespace:"priv"`
	GetEEATransactionCount   func(jsonrpc.Caller) func(context.Context, ethcommon.Address, string, []string) (*hexutil.Uint64, error) `namespace:"priv"`
	GetTransactionCount      func(jsonrpc.Caller) func(context.Context, ethcommon.Address, string) (*hexutil.Uint64, error)           `namespace:"priv"`
}

type BlockNumber int64

func (n BlockNumber) MarshalText() ([]byte, error) {
	switch n {
	case -2:
		return []byte(`pending`), nil
	case -1:
		return []byte(`latest`), nil
	case 0:
		return []byte(`earliest`), nil
	default:
		buf := make([]byte, 2, 10)
		copy(buf, `0x`)
		buf = strconv.AppendUint(buf, uint64(n), 16)
		return buf, nil
	}
}

const (
	PendingBlockNumber  = BlockNumber(-2)
	LatestBlockNumber   = BlockNumber(-1)
	EarliestBlockNumber = BlockNumber(0)
)

// CallMsg contains parameters for contract calls.
type CallMsg struct {
	From     *ethcommon.Address // the sender of the 'transaction'
	To       *ethcommon.Address // the destination contract (nil for contract creation)
	Gas      *uint64            // if 0, the call executes with near-infinite gas
	GasPrice *big.Int           // wei <-> gas exchange ratio
	Value    *big.Int           // amount of wei sent along with the call
	Data     *[]byte            // input data, usually an ABI-encoded contract method invocation
}

func (msg *CallMsg) WithFrom(addr ethcommon.Address) *CallMsg {
	msg.From = &addr
	return msg
}

func (msg *CallMsg) WithTo(addr ethcommon.Address) *CallMsg {
	msg.To = &addr
	return msg
}

func (msg *CallMsg) WithGas(gas uint64) *CallMsg {
	msg.Gas = &gas
	return msg
}

func (msg *CallMsg) WithGasPrice(p *big.Int) *CallMsg {
	msg.GasPrice = p
	return msg
}

func (msg *CallMsg) WithValue(v *big.Int) *CallMsg {
	msg.Value = v
	return msg
}

func (msg *CallMsg) WithData(d []byte) *CallMsg {
	if len(d) != 0 {
		b := make([]byte, len(d))
		copy(b, d)
		msg.Data = &b
	}
	return msg
}

type jsonCallMsg struct {
	From     *ethcommon.Address `json:"from,omitempty"`
	To       *ethcommon.Address `json:"to,omitempty"`
	Gas      *uint64            `json:"gas,omitempty"`
	GasPrice *hexutil.Big       `json:"gasPrice,omitempty"`
	Value    *hexutil.Big       `json:"value,omitempty"`
	Data     *[]byte            `json:"data,omitempty"`
}

func (msg *CallMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonCallMsg{
		From:     msg.From,
		To:       msg.To,
		Gas:      msg.Gas,
		GasPrice: (*hexutil.Big)(msg.GasPrice),
		Value:    (*hexutil.Big)(msg.Value),
		Data:     msg.Data,
	})
}
