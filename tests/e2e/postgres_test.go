// +build e2e

package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/consensys/quorum-key-manager/src/infra/postgres"
	pgclient "github.com/consensys/quorum-key-manager/src/infra/postgres/client"
	"github.com/consensys/quorum-key-manager/tests"
	"github.com/stretchr/testify/require"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/stretchr/testify/suite"
)

func TestPostgres(t *testing.T) {
	s := new(postgresTestSuite)

	s.ctx = context.Background()
	sig := common.NewSignalListener(func(signal os.Signal) {
		s.err = fmt.Errorf("interrupt signal was caught")
		t.FailNow()
	})
	defer sig.Close()

	s.cfg, s.err = tests.NewConfig()
	suite.Run(t, s)
}

type postgresTestSuite struct {
	suite.Suite
	err      error
	ctx      context.Context
	pgClient postgres.Client
	cfg      *tests.Config
}

type FakeType struct {
	tableName struct{} `pg:"fake_types"`

	ID    int
	Value string
}

func (s *postgresTestSuite) SetupSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	s.pgClient, s.err = pgclient.NewClient(s.cfg.Postgres)
	if s.err != nil {
		s.T().Error(s.err)
	}

	s.err = s.pgClient.CreateTable(s.ctx, &FakeType{})
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func (s *postgresTestSuite) TearDownSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	s.err = s.pgClient.DropTable(s.ctx, &FakeType{})
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func (s *postgresTestSuite) TestInsert() {
	testCase := "insert simple"
	s.Run(testCase, func() {
		data := FakeType{
			ID:    1,
			Value: testCase,
		}

		err := s.pgClient.Insert(s.ctx, &data)
		require.NoError(s.T(), err)
	})
}
