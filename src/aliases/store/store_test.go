package aliasstore_test

import (
	"context"
	"fmt"
	"testing"

	aliases "github.com/consensys/quorum-key-manager/src/aliases"
	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
	"github.com/go-pg/pg/v10"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
)

type closeFunc func()

func startAndConnectDB(t *testing.T) (*pg.DB, closeFunc) {
	t.Helper()
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("can not connect to docker: %s", err)
	}

	pwd := "postgres"
	ctr, err := pool.Run("postgres", "13", []string{fmt.Sprintf("POSTGRES_PASSWORD=%s", pwd)})
	if err != nil {
		t.Fatalf("can not run container: %s", err)
	}
	closeFn := func() {
		err := pool.Purge(ctr)
		if err != nil {
			t.Fatal("can not purge container:", err)
		}
	}

	var db *pg.DB
	opt := pg.Options{
		Password: pwd,
		Addr:     fmt.Sprintf("localhost:%s", ctr.GetPort("5432/tcp")),
	}
	db = pg.Connect(&opt)
	err = pool.Retry(func(db *pg.DB) func() error {
		return func() error {
			_, err := db.Exec("SELECT 1;")
			return err
		}
	}(db))
	if err != nil {
		t.Fatalf("can not connect")
	}
	return db, closeFn
}

func TestCreate(t *testing.T) {
	db, closeFn := startAndConnectDB(t)
	defer closeFn()
	s := aliasstore.New(db)
	ctx := context.TODO()
	alias := aliases.Alias{
		RegistryID: "JPM",
		ID:         "Goldman Sachs",
	}
	db.ModelContext(ctx, &alias).CreateTable(nil)

	err := s.CreateAlias(ctx, alias.RegistryID, alias.ID)
	require.NoError(t, err)

	a, err := s.GetAlias(ctx, alias.RegistryID, alias.ID)

	require.NoError(t, err)
	require.Equal(t, a.RegistryID, alias.RegistryID)
	require.Equal(t, a.ID, alias.ID)
}
