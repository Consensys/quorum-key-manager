// +build e2e

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/tests"
	"github.com/stretchr/testify/suite"
)

type healthzTestSuite struct {
	suite.Suite
	err          error
	ctx          context.Context
	healthClient *http.Client
	cfg          *tests.Config
	mainAccount  *types.Eth1AccountResponse
}

func (s *healthzTestSuite) SetupSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	s.healthClient = &http.Client{}
}

func (s *healthzTestSuite) TearDownSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func TestHealthz(t *testing.T) {
	s := new(healthzTestSuite)

	s.ctx = context.Background()
	sig := common.NewSignalListener(func(signal os.Signal) {
		s.err = fmt.Errorf("interrupt signal was caught")
		t.FailNow()
	})
	defer sig.Close()

	s.cfg, s.err = tests.NewConfig()
	suite.Run(t, s)
}

func (s *healthzTestSuite) TestLiveness() {
	s.Run("should validate liveness endpoint", func() {
		
	})
}

func (s *healthzTestSuite) TestReadiness() {
	s.Run("should validate readiness endpoint", func() {
		
	})
}

func CheckLiveness(ctx context.Context, client *http.Client, healthURL string) (bool, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/live", healthURL), nil)

	res, err := client.Do(req)
	if err != nil {
		return false, err
	}

	return res.StatusCode == http.StatusOK, nil
}

func CheckReadiness(ctx context.Context, client *http.Client, healthURL string) (bool, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/ready", healthURL), nil)

	res, err := client.Do(req)
	if err != nil {
		return false, err
	}

	return res.StatusCode == http.StatusOK, nil
}
