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

func fakeAlias() aliases.Alias {
	return aliases.Alias{
		RegistryID: "JPM",
		ID:         "Goldman Sachs",
		Kind:       aliases.AliasKindArray,
		Value:      `["ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=","2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="]`,
	}
}
func TestCreate(t *testing.T) {
	t.Parallel()
	db, closeFn := startAndConnectDB(t)
	defer closeFn()
	s := aliasstore.New(db)

	in := fakeAlias()
	ctx := context.Background()
	db.ModelContext(ctx, &in).CreateTable(nil)

	err := s.CreateAlias(ctx, in)
	require.NoError(t, err)
}

func TestGet(t *testing.T) {
	t.Parallel()
	db, closeFn := startAndConnectDB(t)
	defer closeFn()
	s := aliasstore.New(db)

	in := fakeAlias()
	ctx := context.Background()
	db.ModelContext(ctx, &in).CreateTable(nil)
	t.Run("non existing alias", func(t *testing.T) {
		_, err := s.GetAlias(ctx, in.RegistryID, in.ID)
		require.NotNil(t, err)
	})

	t.Run("just created alias", func(t *testing.T) {
		err := s.CreateAlias(ctx, in)
		require.NoError(t, err)

		got, err := s.GetAlias(ctx, in.RegistryID, in.ID)
		require.NoError(t, err)
		require.Equal(t, &in, got)
	})
}

}
