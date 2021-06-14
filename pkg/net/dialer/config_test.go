package dialer

import (
	"testing"
	"time"

	"github.com/consensysquorum/quorum-key-manager/pkg/json"
	"github.com/stretchr/testify/assert"
)

func TestConfigSetDefault(t *testing.T) {
	tests := []struct {
		desc string

		cfg         *Config
		expectedCfg *Config
	}{
		{
			desc: "empty",
			cfg:  &Config{},
			expectedCfg: &Config{
				Timeout:   &json.Duration{Duration: 30 * time.Second},
				KeepAlive: &json.Duration{Duration: 30 * time.Second},
			},
		},
		{
			desc: "non empty",
			cfg: &Config{
				Timeout:   &json.Duration{Duration: 60 * time.Second},
				KeepAlive: &json.Duration{Duration: 60 * time.Second},
			},
			expectedCfg: &Config{
				Timeout:   &json.Duration{Duration: 60 * time.Second},
				KeepAlive: &json.Duration{Duration: 60 * time.Second},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.cfg.SetDefault()
			assert.Equal(t, *tt.cfg.KeepAlive, *tt.expectedCfg.KeepAlive, "Timeout should match")
			assert.Equal(t, *tt.cfg.KeepAlive, *tt.expectedCfg.KeepAlive, "KeepAlive should match")
		})
	}
}
