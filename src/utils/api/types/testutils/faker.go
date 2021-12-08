package testutils

import (
	"github.com/consensys/quorum-key-manager/src/stores/api/types/testutils"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/consensys/quorum-key-manager/src/utils/api/types"
)

func FakeECRecoverRequest() *types.ECRecoverRequest {
	return &types.ECRecoverRequest{
		Data:      hexutil.MustDecode("0xfeee"),
		Signature: hexutil.MustDecode("0x3399aeb23d6564b3a0b220447e9f1bb2057ffb82cfb766147620aa6bc84938e26941e7583d6460fea405d99da897e88cab07a7fd0991c6c2163645c45d25e4b201"),
	}
}

func FakeVerifyTypedDataPayloadRequest() *types.VerifyTypedDataRequest {
	return &types.VerifyTypedDataRequest{
		TypedData: *testutils.FakeSignTypedDataRequest(),
		Signature: hexutil.MustDecode("0x3399aeb23d6564b3a0b220447e9f1bb2057ffb82cfb766147620aa6bc84938e26941e7583d6460fea405d99da897e88cab07a7fd0991c6c2163645c45d25e4b201"),
		Address:   common.HexToAddress("0x5Cc634233E4a454d47aACd9fC68801482Fb02610"),
	}
}

func FakeVerifyRequest() *types.VerifyRequest {
	return &types.VerifyRequest{
		Data:      hexutil.MustDecode("0xfeee"),
		Signature: hexutil.MustDecode("0x3399aeb23d6564b3a0b220447e9f1bb2057ffb82cfb766147620aa6bc84938e26941e7583d6460fea405d99da897e88cab07a7fd0991c6c2163645c45d25e4b201"),
		Address:   common.HexToAddress("0x5Cc634233E4a454d47aACd9fC68801482Fb02610"),
	}
}
