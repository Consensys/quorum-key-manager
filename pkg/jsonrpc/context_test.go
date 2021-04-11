package jsonrpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	req := RequestFromContext(context.Background())
	assert.Nil(t, req)

	req = new(Request)
	ctx := WithRequest(context.Background(), req)
	assert.Equal(t, req, RequestFromContext(ctx))
}
