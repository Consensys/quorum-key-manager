package node

import (
	"bytes"
	"context"
	"io/ioutil"
	"math/big"
	"net/http"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEthService(t *testing.T) {
	srv := new(EthService)
	err := jsonrpc.Provide(srv)
	require.NoError(t, err, "Provide must not error")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transport := testutils.NewMockRoundTripper(ctrl)
	client := jsonrpc.NewClient(&http.Client{Transport: transport})
	req, _ := http.NewRequest(http.MethodPost, "www.example.com", nil)

	// Empty ID and version client
	cllr := jsonrpc.NewCaller(jsonrpc.WithVersion("2.0")(client), jsonrpc.NewRequest(req))

	// GasPrice
	t.Run("GasPrice", func(t *testing.T) {
		m := testutils.RequestMatcher(
			t,
			"www.example.com",
			[]byte(`{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":null}`),
		)
		respBody := []byte(`{"jsonrpc": "2.0","result":"0x3e8"}`)
		transport.EXPECT().RoundTrip(m).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		}, nil)

		gasPrice, err := srv.GasPrice(cllr)(context.Background())
		require.NoError(t, err, "Must not error")
		assert.Equal(t, "0x3e8", gasPrice.String(), "Result should be valid")
	})

	// GetTransactionCount on pending block
	t.Run("GetTransactionCount on pending", func(t *testing.T) {
		m := testutils.RequestMatcher(
			t,
			"www.example.com",
			[]byte(`{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["0xc94770007dda54cf92009bff0de90c06f603a09f","latest"],"id":null}`),
		)
		respBody := []byte(`{"jsonrpc": "2.0","result":"0xf"}`)
		transport.EXPECT().RoundTrip(m).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		}, nil)

		count, err := srv.GetTransactionCount(cllr)(context.Background(), ethcommon.HexToAddress("0xc94770007dda54cF92009BFF0dE90c06F603a09f"), LatestBlockNumber)
		require.NoError(t, err, "Must not error")
		assert.Equal(t, uint64(15), uint64(*count), "Result should be valid")
	})

	t.Run("GetTransactionCount on numbered block", func(t *testing.T) {
		m := testutils.RequestMatcher(
			t,
			"www.example.com",
			[]byte(`{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["0xc94770007dda54cf92009bff0de90c06f603a09f","0xa"],"id":null}`),
		)
		respBody := []byte(`{"jsonrpc": "2.0","result":"0xf"}`)
		transport.EXPECT().RoundTrip(m).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		}, nil)

		count, err := srv.GetTransactionCount(cllr)(context.Background(), ethcommon.HexToAddress("0xc94770007dda54cF92009BFF0dE90c06F603a09f"), BlockNumber(10))
		require.NoError(t, err, "Must not error")
		assert.Equal(t, uint64(15), uint64(*count), "Result should be valid")
	})

	t.Run("SendRawTransaction", func(t *testing.T) {
		m := testutils.RequestMatcher(
			t,
			"www.example.com",
			[]byte(`{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0xf869018203e882520894f17f52151ebef6c7334fad080c5704d77216b732881bc16d674ec80000801ba02da1c48b670996dcb1f447ef9ef00b33033c48a4fe938f420bec3e56bfd24071a062e0aa78a81bf0290afbc3a9d8e9a068e6d74caa66c5e0fa8a46deaae96b0833"],"id":null}`),
		)
		respBody := []byte(`{"jsonrpc": "2.0","result":"0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331a"}`)
		transport.EXPECT().RoundTrip(m).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		}, nil)

		hash, err := srv.SendRawTransaction(cllr)(context.Background(), ethcommon.FromHex("0xf869018203e882520894f17f52151ebef6c7334fad080c5704d77216b732881bc16d674ec80000801ba02da1c48b670996dcb1f447ef9ef00b33033c48a4fe938f420bec3e56bfd24071a062e0aa78a81bf0290afbc3a9d8e9a068e6d74caa66c5e0fa8a46deaae96b0833"))
		require.NoError(t, err, "Must not error")
		assert.Equal(t, "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331a", hash.String(), "Result should be valid")
	})

	t.Run("EstimateGas", func(t *testing.T) {
		m := testutils.RequestMatcher(
			t,
			"www.example.com",
			[]byte(`{"jsonrpc":"2.0","method":"eth_estimateGas","params":[{"from":"0xfe3b557e8fb62b89f4916b721be55ceb828dbd73","to":"0x44aa93095d6749a706051658b970b941c72c1d53","value":"0x1"}],"id":null}`),
		)
		respBody := []byte(`{"jsonrpc": "2.0","result":"0x5208"}`)
		transport.EXPECT().RoundTrip(m).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		}, nil)

		msg := (&CallMsg{}).
			WithTo(ethcommon.HexToAddress("0x44Aa93095D6749A706051658B970b941c72c1D53")).
			WithFrom(ethcommon.HexToAddress("0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73")).
			WithValue(big.NewInt(1))
		gas, err := srv.EstimateGas(cllr)(context.Background(), msg)
		require.NoError(t, err, "Must not error")
		assert.Equal(t, uint64(21000), uint64(*gas), "Result should be valid")
	})
}

func TestEEAService(t *testing.T) {
	srv := new(EEAService)
	err := jsonrpc.Provide(srv)
	require.NoError(t, err, "Provide must not error")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transport := testutils.NewMockRoundTripper(ctrl)
	client := jsonrpc.NewClient(&http.Client{Transport: transport})
	req, _ := http.NewRequest(http.MethodPost, "www.example.com", nil)

	// Empty ID and version client
	cllr := jsonrpc.NewCaller(jsonrpc.WithVersion("2.0")(client), jsonrpc.NewRequest(req))

	t.Run("SendRawTransaction", func(t *testing.T) {
		m := testutils.RequestMatcher(
			t,
			"www.example.com",
			[]byte(`{"jsonrpc":"2.0","method":"eea_sendRawTransaction","params":["0xf869018203e882520894f17f52151ebef6c7334fad080c5704d77216b732881bc16d674ec80000801ba02da1c48b670996dcb1f447ef9ef00b33033c48a4fe938f420bec3e56bfd24071a062e0aa78a81bf0290afbc3a9d8e9a068e6d74caa66c5e0fa8a46deaae96b0833"],"id":null}`),
		)
		respBody := []byte(`{"jsonrpc": "2.0","result":"0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331a"}`)
		transport.EXPECT().RoundTrip(m).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		}, nil)

		hash, err := srv.SendRawTransaction(cllr)(context.Background(), ethcommon.FromHex("0xf869018203e882520894f17f52151ebef6c7334fad080c5704d77216b732881bc16d674ec80000801ba02da1c48b670996dcb1f447ef9ef00b33033c48a4fe938f420bec3e56bfd24071a062e0aa78a81bf0290afbc3a9d8e9a068e6d74caa66c5e0fa8a46deaae96b0833"))
		require.NoError(t, err, "Must not error")
		assert.Equal(t, "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331a", hash.String(), "Result should be valid")
	})
}

func TestPrivService(t *testing.T) {
	srv := new(PrivService)
	err := jsonrpc.Provide(srv)
	require.NoError(t, err, "Provide must not error")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transport := testutils.NewMockRoundTripper(ctrl)
	client := jsonrpc.NewClient(&http.Client{Transport: transport})
	req, _ := http.NewRequest(http.MethodPost, "www.example.com", nil)

	// Empty ID and version client
	cllr := jsonrpc.NewCaller(jsonrpc.WithVersion("2.0")(client), jsonrpc.NewRequest(req))

	t.Run("DistributeRawTransaction", func(t *testing.T) {
		m := testutils.RequestMatcher(
			t,
			"www.example.com",
			[]byte(`{"jsonrpc":"2.0","method":"priv_distributeRawTransaction","params":["0xf869018203e882520894f17f52151ebef6c7334fad080c5704d77216b732881bc16d674ec80000801ba02da1c48b670996dcb1f447ef9ef00b33033c48a4fe938f420bec3e56bfd24071a062e0aa78a81bf0290afbc3a9d8e9a068e6d74caa66c5e0fa8a46deaae96b0833"],"id":null}`),
		)
		respBody := []byte(`{"jsonrpc": "2.0","result":"0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331a"}`)
		transport.EXPECT().RoundTrip(m).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		}, nil)

		enclaveKey, err := srv.DistributeRawTransaction(cllr)(context.Background(), ethcommon.FromHex("0xf869018203e882520894f17f52151ebef6c7334fad080c5704d77216b732881bc16d674ec80000801ba02da1c48b670996dcb1f447ef9ef00b33033c48a4fe938f420bec3e56bfd24071a062e0aa78a81bf0290afbc3a9d8e9a068e6d74caa66c5e0fa8a46deaae96b0833"))
		require.NoError(t, err, "Must not error")
		assert.Equal(t, "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331a", enclaveKey.String(), "Result should be valid")
	})

	t.Run("GetEEATransactionCount", func(t *testing.T) {
		m := testutils.RequestMatcher(
			t,
			"www.example.com",
			[]byte(`{"jsonrpc":"2.0","method":"priv_getEEATransactionCount","params":["0xc94770007dda54cf92009bff0de90c06f603a09f","GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=",["KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=","eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="]],"id":null}`),
		)
		respBody := []byte(`{"jsonrpc": "2.0","result":"0xf"}`)
		transport.EXPECT().RoundTrip(m).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		}, nil)

		count, err := srv.GetEEATransactionCount(cllr)(context.Background(), ethcommon.HexToAddress("0xc94770007dda54cF92009BFF0dE90c06F603a09f"), "GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=", []string{"KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=", "eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="})
		require.NoError(t, err, "Must not error")
		assert.Equal(t, uint64(15), uint64(*count), "Result should be valid")
	})

	t.Run("GetTransactionCount", func(t *testing.T) {
		m := testutils.RequestMatcher(
			t,
			"www.example.com",
			[]byte(`{"jsonrpc":"2.0","method":"priv_getTransactionCount","params":["0xc94770007dda54cf92009bff0de90c06f603a09f","kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M="],"id":null}`),
		)
		respBody := []byte(`{"jsonrpc": "2.0","result":"0xf"}`)
		transport.EXPECT().RoundTrip(m).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		}, nil)

		count, err := srv.GetTransactionCount(cllr)(context.Background(), ethcommon.HexToAddress("0xc94770007dda54cF92009BFF0dE90c06F603a09f"), "kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M=")
		require.NoError(t, err, "Must not error")
		assert.Equal(t, uint64(15), uint64(*count), "Resut should be valid")
	})
}
