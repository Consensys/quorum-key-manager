// +build integration

package integrationtests

import (
	"context"
	"os"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/integration-tests/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type storeTestSuite struct {
	suite.Suite
	env             *IntegrationEnvironment
	err             error
}

func (s *storeTestSuite) SetupSuite() {
	err := StartEnvironment(s.env.ctx, s.env)
	if err != nil {
		s.T().Error(err)
		return
	}

	s.env.logger.Info("setup test suite has completed")
}

func (s *storeTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func TestKeyManagerStore(t *testing.T) {
	s := new(storeTestSuite)
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

func (s *storeTestSuite) TestKeyManagerStore_Template() {
	s.T().Run("should success", func(t *testing.T) {
		assert.True(t, true)
	})
}
