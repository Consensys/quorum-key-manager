// +build acceptance

package integrationtests

import (
	"context"
	"os"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests/utils"
	"github.com/stretchr/testify/suite"
)

type keyManagerTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
	err error
}

func (s *keyManagerTestSuite) SetupSuite() {
	err := StartEnvironment(s.env.ctx, s.env)
	if err != nil {
		s.T().Error(err)
		return
	}

	s.env.logger.Info("setup test suite has completed")
}

func (s *keyManagerTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func TestKeyManager(t *testing.T) {
	s := new(keyManagerTestSuite)
	ctx, cancel := context.WithCancel(context.Background())

	var err error
	s.env, err = NewIntegrationEnvironment(ctx)
	if err != nil {
		t.Error(err.Error())
		return
	}

	sig := utils.NewSignalListener(func(signal os.Signal) {
		cancel()
	})
	defer sig.Close()

	suite.Run(t, s)
}

func (s *keyManagerTestSuite) TestKeyManager_Secrets() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(secretsTestSuite)
	testSuite.env = s.env
	testSuite.baseURL = s.env.baseURL

	suite.Run(s.T(), testSuite)
}
