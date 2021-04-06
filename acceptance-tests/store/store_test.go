// +build acceptance

package store

import (
	"context"
	"github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets/hashicorp"
	"os"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/stretchr/testify/suite"
)

type storeTestSuite struct {
	suite.Suite
	env *integrationtests.IntegrationEnvironment
	err error
}

func (s *storeTestSuite) SetupSuite() {
	err := integrationtests.StartEnvironment(s.env.ctx, s.env)
	if err != nil {
		s.T().Error(err)
		return
	}

	s.env.Logger.Info("setup test suite has completed")
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
	s.env, err = integrationtests.NewIntegrationEnvironment(ctx)
	if err != nil {
		t.Error(err.Error())
		return
	}

	sig := common.NewSignalListener(func(signal os.Signal) {
		cancel()
	})
	defer sig.Close()

	suite.Run(t, s)
}

func (s *storeTestSuite) TestKeyManagerStore_HashicorpSecret() {
	if s.err != nil {
		s.env.Logger.Warn("skipping test...")
		return
	}

	store := hashicorp.New(s.env.HashicorpClient, "orchestrate-hashicorp-vault-plugin")

	testSuite := new(hashicorpSecretTestSuite)
	testSuite.env = s.env
	testSuite.store = store
	suite.Run(s.T(), testSuite)
}
