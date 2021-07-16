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
		err = pool.Purge(ctr)
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
			_, err = db.Exec("SELECT 1;")
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
func TestCreateAlias(t *testing.T) {
	t.Parallel()
	db, closeFn := startAndConnectDB(t)
	defer closeFn()
	s := aliasstore.New(db)

	in := fakeAlias()
	ctx := context.Background()
	err := db.ModelContext(ctx, &in).CreateTable(nil)
	require.NoError(t, err)

	err = s.CreateAlias(ctx, in)
	require.NoError(t, err)
}

func TestGetAlias(t *testing.T) {
	t.Parallel()
	db, closeFn := startAndConnectDB(t)
	defer closeFn()
	s := aliasstore.New(db)

	in := fakeAlias()
	ctx := context.Background()
	err := db.ModelContext(ctx, &in).CreateTable(nil)
	require.NoError(t, err)
	t.Run("non existing alias", func(t *testing.T) {
		_, err := s.GetAlias(ctx, in.RegistryID, in.ID)
		require.Error(t, err)
	})

	t.Run("just created alias", func(t *testing.T) {
		err := s.CreateAlias(ctx, in)
		require.NoError(t, err)

		got, err := s.GetAlias(ctx, in.RegistryID, in.ID)
		require.NoError(t, err)
		require.Equal(t, &in, got)
	})
}

func TestUpdateAlias(t *testing.T) {
	t.Parallel()
	db, closeFn := startAndConnectDB(t)
	defer closeFn()
	s := aliasstore.New(db)

	in := fakeAlias()
	ctx := context.Background()
	err := db.ModelContext(ctx, &in).CreateTable(nil)
	require.NoError(t, err)
	t.Run("non existing alias", func(t *testing.T) {
		err := s.UpdateAlias(ctx, in)
		require.Error(t, err)
	})

	t.Run("just created alias", func(t *testing.T) {
		err := s.CreateAlias(ctx, in)
		require.NoError(t, err)

		updated := in
		updated.Value = `["SOAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=","3T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="]`

		err = s.UpdateAlias(ctx, updated)
		require.NoError(t, err)

		got, err := s.GetAlias(ctx, in.RegistryID, in.ID)
		require.NoError(t, err)
		require.Equal(t, &updated, got)
	})
}

func TestDeleteAlias(t *testing.T) {
	t.Parallel()
	db, closeFn := startAndConnectDB(t)
	defer closeFn()
	s := aliasstore.New(db)

	in := fakeAlias()
	ctx := context.Background()
	err := db.ModelContext(ctx, &in).CreateTable(nil)
	require.NoError(t, err)
	t.Run("non existing alias", func(t *testing.T) {
		err := s.DeleteAlias(ctx, in.RegistryID, in.ID)
		require.Error(t, err)
	})

	t.Run("just created alias", func(t *testing.T) {
		err := s.CreateAlias(ctx, in)
		require.NoError(t, err)

		err = s.DeleteAlias(ctx, in.RegistryID, in.ID)
		require.NoError(t, err)

		_, err = s.GetAlias(ctx, in.RegistryID, in.ID)
		require.Error(t, err)
	})
}

func TestListAlias(t *testing.T) {
	t.Parallel()
	db, closeFn := startAndConnectDB(t)
	defer closeFn()
	s := aliasstore.New(db)

	in := fakeAlias()
	ctx := context.Background()
	err := db.ModelContext(ctx, &in).CreateTable(nil)
	require.NoError(t, err)
	t.Run("non existing alias", func(t *testing.T) {
		als, err := s.ListAliases(ctx, in.RegistryID)
		require.NoError(t, err)
		require.Len(t, als, 0)
	})

	t.Run("just created alias", func(t *testing.T) {
		err := s.CreateAlias(ctx, in)
		require.NoError(t, err)

		newAlias := in
		newAlias.ID = `Crédit Mutuel`
		newAlias.Kind = aliases.AliasKindString
		newAlias.Value = `SOAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=`
		err = s.CreateAlias(ctx, newAlias)
		require.NoError(t, err)

		als, err := s.ListAliases(ctx, in.RegistryID)
		require.NoError(t, err)
		require.NotEmpty(t, als)
		require.Len(t, als, 2)
		require.Equal(t, als[0].ID, in.ID)
		require.Equal(t, als[1].ID, newAlias.ID)
	})
}

func TestDeleteRegistry(t *testing.T) {
	t.Parallel()
	db, closeFn := startAndConnectDB(t)
	defer closeFn()
	s := aliasstore.New(db)

	in := fakeAlias()
	ctx := context.Background()
	err := db.ModelContext(ctx, &in).CreateTable(nil)
	require.NoError(t, err)
	t.Run("non existing alias", func(t *testing.T) {
		err := s.DeleteAlias(ctx, in.RegistryID, in.ID)
		require.Error(t, err)
	})

	t.Run("just created alias", func(t *testing.T) {
		err := s.CreateAlias(ctx, in)
		require.NoError(t, err)

		err = s.DeleteAlias(ctx, in.RegistryID, in.ID)
		require.NoError(t, err)

		_, err = s.GetAlias(ctx, in.RegistryID, in.ID)
		require.Error(t, err)
	})
}
