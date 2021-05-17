package proxynode

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	sess := SessionFromContext(context.Background())
	assert.Nil(t, sess)

	s := &session{}
	ctx := WithSession(context.Background(), s)
	assert.Equal(t, s, SessionFromContext(ctx))
}
