// +build e2e

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type healthzTestSuite struct {
	suite.Suite
	err         error
	mainAccount *types.EthAccountResponse
	env         *Environment
	client      *http.Client
}

func TestHealthz(t *testing.T) {
	s := new(healthzTestSuite)

	sig := common.NewSignalListener(func(signal os.Signal) {
		s.err = fmt.Errorf("interrupt signal was caught")
		t.FailNow()
	})
	defer sig.Close()

	env, err := NewEnvironment()
	require.NoError(t, err)
	s.env = env

	s.client = &http.Client{Transport: NewTestHttpTransport("", "", nil)}

	suite.Run(t, s)
}

func (s *healthzTestSuite) TestLiveness() {

	s.Run("should validate liveness endpoint", func() {
		isLive, err := s.checkLiveness(s.env.ctx)
		require.NoError(s.T(), err)

		assert.True(s.T(), isLive)
	})
}

func (s *healthzTestSuite) TestReadiness() {
	s.Run("should validate readiness endpoint", func() {
		isReady, err := s.checkReadiness(s.env.ctx)
		require.NoError(s.T(), err)

		assert.True(s.T(), isReady)
	})
}

func (s *healthzTestSuite) checkLiveness(ctx context.Context) (bool, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/live", s.env.cfg.HealthKeyManagerURL), nil)

	res, err := s.client.Do(req)
	if err != nil {
		return false, err
	}

	return res.StatusCode == http.StatusOK, nil
}

func (s *healthzTestSuite) checkReadiness(ctx context.Context) (bool, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/ready", s.env.cfg.HealthKeyManagerURL), nil)

	res, err := s.client.Do(req)
	if err != nil {
		return false, err
	}

	return res.StatusCode == http.StatusOK, nil
}
