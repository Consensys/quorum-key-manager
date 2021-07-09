package aliasstore_test

import (
	"context"
	"fmt"
	"testing"

	aliases "github.com/consensys/quorum-key-manager/src/aliases"
	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
	"github.com/go-pg/pg/v10"
	"github.com/ory/dockertest"
)

func TestCreate(t *testing.T) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("can not connect to docker: %s", err)
	}

	pwd := "postgres"
	ctr, err := pool.Run("postgres", "10", []string{fmt.Sprintf("POSTGRES_PASSWORD=%s", pwd)})
	if err != nil {
		t.Fatalf("can not run container: %s", err)
	}
	defer func() {
		err := pool.Purge(ctr)
		if err != nil {
			t.Fatal("can not purge container:", err)
		}
	}()

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
	s := aliasstore.New(db)
	ctx := context.TODO()
	alias := aliases.Alias{
		RegistryID: "JPM",
		ID:         "Goldman Sachs",
	}
	db.ModelContext(ctx, &alias).CreateTable(nil)

	err = s.CreateAlias(ctx, alias.RegistryID, alias.ID)
	if err != nil {
		t.Fatal("can not create after connect:", err)
	}

	a, err := s.GetAlias(ctx, alias.RegistryID, alias.ID)
	t.Log("alias:", a)
	if err != nil {
		t.Fatal("can not select after connect:", err)
	}
	if a.RegistryID != alias.RegistryID || a.ID != alias.ID {
		t.Fatal("did not get same object", err)
	}
}
