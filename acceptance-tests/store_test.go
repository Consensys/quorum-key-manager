// +build acceptance

package integrationtests

import (
	"context"
	hashicorpkey "github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys/hashicorp"
	"os"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/stretchr/testify/suite"
)

type storeTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
	err error
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

	sig := common.NewSignalListener(func(signal os.Signal) {
		cancel()
	})
	defer sig.Close()

	suite.Run(t, s)
}

/*
func (s *storeTestSuite) TestKeyManagerStore_HashicorpSecret() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	store := hashicorpsecret.New(s.env.hashicorpClient, "secret")

	testSuite := new(hashicorpSecretTestSuite)
	testSuite.env = s.env
	testSuite.store = store
	suite.Run(s.T(), testSuite)
}

func (s *storeTestSuite) TestKeyManagerStore_AKVSecret() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	store := akv.New(s.env.akvClient)

	testSuite := new(akvSecretTestSuite)
	testSuite.env = s.env
	testSuite.store = store
	suite.Run(s.T(), testSuite)
}
*/

func (s *storeTestSuite) TestKeyManagerStore_HashicorpKey() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	store := hashicorpkey.New(s.env.hashicorpClient, "orchestrate")

	testSuite := new(hashicorpKeyTestSuite)
	testSuite.env = s.env
	testSuite.store = store
	suite.Run(s.T(), testSuite)
}
