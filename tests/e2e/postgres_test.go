// +build e2e

package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
	pgclient "github.com/consensys/quorum-key-manager/src/infra/postgres/client"
	"github.com/consensys/quorum-key-manager/tests"
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

func (s *postgresTestSuite) TestSelect() {
	testCase := "non-existing select simple"
	s.Run(testCase, func() {
		data := FakeType{
			ID:    2,
			Value: testCase,
		}

		got := data
		err := s.pgClient.Select(s.ctx, &got)
		require.Error(s.T(), err)
		require.True(s.T(), errors.IsNotFoundError(err))
	})

	testCase = "select simple"
	s.Run(testCase, func() {
		data := FakeType{
			ID:    21,
			Value: testCase,
		}

		err := s.pgClient.Insert(s.ctx, &data)
		require.NoError(s.T(), err)

		got := data
		err = s.pgClient.Select(s.ctx, &got)
		require.NoError(s.T(), err)
		require.Equal(s.T(), data.ID, got.ID)
		require.Equal(s.T(), data.Value, got.Value)
	})
}

func (s *postgresTestSuite) TestUpdate() {
	testCase := "update simple"
	s.Run(testCase, func() {
		data := FakeType{
			ID:    3,
			Value: testCase,
		}

		err := s.pgClient.Insert(s.ctx, &data)
		require.NoError(s.T(), err)

		updated := data
		updated.Value = "update simple: updated"

		err = s.pgClient.Update(s.ctx, &updated)
		require.NoError(s.T(), err)

		got := FakeType{ID: 3, Value: "wrong"}
		err = s.pgClient.Select(s.ctx, &got)
		require.Equal(s.T(), updated.ID, got.ID)
		require.NotEqual(s.T(), data.Value, got.Value)
		require.Equal(s.T(), updated.Value, got.Value)
	})
}

func (s *postgresTestSuite) TestDelete() {
	testCase := "delete simple"
	s.Run(testCase, func() {
		data := FakeType{
			ID:    4,
			Value: testCase,
		}

		err := s.pgClient.Insert(s.ctx, &data)
		require.NoError(s.T(), err)

		toDelete := data
		err = s.pgClient.Delete(s.ctx, &toDelete)
		require.NoError(s.T(), err)

		got := data
		err = s.pgClient.Select(s.ctx, &got)
		require.Error(s.T(), err)
	})
}
