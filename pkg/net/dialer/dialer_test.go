package dialer

import "testing"

func TestDialer(t *testing.T) {
	_ = New(new(Config).SetDefault())
}
