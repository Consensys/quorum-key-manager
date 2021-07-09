package nodes

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	nodesapi "github.com/consensys/quorum-key-manager/src/nodes/api"
	nodesmanager "github.com/consensys/quorum-key-manager/src/nodes/manager"
	storesmanager "github.com/consensys/quorum-key-manager/src/stores/manager"
)

func RegisterService(a *app.App, logger log.Logger) error {
	// Load manifests service
	m := new(manifestsmanager.Manager)
	err := a.Service(m)
	if err != nil {
		return err
	}

	// Load stores service
	s := new(storesmanager.Manager)
	err = a.Service(s)
	if err != nil {
		return err
	}

	// Create and register nodes service
	nodes := nodesmanager.New(*s, *m, logger)
	err = a.RegisterService(nodes)
	if err != nil {
		return err
	}

	// Create and register nodes API
	nodesapi.New(nodes).Register(a.Router())

	return nil
}
