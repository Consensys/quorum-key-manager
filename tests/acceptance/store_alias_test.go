// +build e2e

package aliasstore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"

	"github.com/consensys/quorum-key-manager/src/aliases"
	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
	"github.com/consensys/quorum-key-manager/src/infra/postgres/client"
)

type closeFunc func()

func startAndConnectDB(ctx context.Context, t *testing.T) (postgres.Client, closeFunc) {
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

	pgClient, err := client.NewClient(&client.Config{
		Password: pwd,
		Host:     "localhost",
		Port:     ctr.GetPort("5432/tcp"),
	})
	err = pool.Retry(func(cl postgres.Client) func() error {
		return func() error {
			err = cl.Ping(ctx)
			return err
		}
	}(pgClient))
	if err != nil {
		t.Fatalf("can not connect")
	}

	return pgClient, closeFn
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
	ctx := context.Background()
	db, closeFn := startAndConnectDB(ctx, t)
	defer closeFn()

	s := aliasstore.New(db)

	in := fakeAlias()
	err := db.CreateTable(ctx, &in)
	require.NoError(t, err)

	err = s.CreateAlias(ctx, in)
	require.NoError(t, err)
}

func TestGetAlias(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	db, closeFn := startAndConnectDB(ctx, t)
	defer closeFn()
	s := aliasstore.New(db)

	in := fakeAlias()
	err := db.CreateTable(ctx, &in)
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
	ctx := context.Background()
	db, closeFn := startAndConnectDB(ctx, t)
	defer closeFn()
	s := aliasstore.New(db)

	in := fakeAlias()
	err := db.CreateTable(ctx, &in)
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
	ctx := context.Background()
	db, closeFn := startAndConnectDB(ctx, t)
	defer closeFn()
	s := aliasstore.New(db)

	in := fakeAlias()
	err := db.CreateTable(ctx, &in)
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
	ctx := context.Background()
	db, closeFn := startAndConnectDB(ctx, t)
	defer closeFn()
	s := aliasstore.New(db)

	in := fakeAlias()
	err := db.CreateTable(ctx, &in)
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
