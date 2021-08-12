package eth1

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/ethereum/go-ethereum/signer/core"
)

func (c Connector) Verify(ctx context.Context, addr string, data, sig []byte) error {
	recoveredAddress, err := c.ECRecover(ctx, data, sig)
	if err != nil {
		return err
	}

	if addr != recoveredAddress {
		errMessage := "failed to verify signature: recovered address does not match the expected one or payload is malformed"
		c.logger.WithError(err).Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	c.logger.Debug("data verified successfully")
	return nil
}

func (c Connector) VerifyTypedData(ctx context.Context, addr string, typedData *core.TypedData, sig []byte) error {
	encodedData, err := getEIP712EncodedData(typedData)
	if err != nil {
		errMessage := "failed to generate EIP-712 encoded data"
		c.logger.WithError(err).Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	c.logger.Debug("typed data verified successfully")
	return c.Verify(ctx, addr, []byte(encodedData), sig)
}
