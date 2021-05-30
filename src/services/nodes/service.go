package nodes

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/app"
	manifestsmanager "github.com/ConsenSysQuorum/quorum-key-manager/src/services/manifests/manager"
	nodesapi "github.com/ConsenSysQuorum/quorum-key-manager/src/services/nodes/api"
	nodesmanager "github.com/ConsenSysQuorum/quorum-key-manager/src/services/nodes/manager"
	storesmanager "github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/manager"
)

func RegisterService(a *app.App) error {
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
	nodes := nodesmanager.New(*s, *m)
	err = a.RegisterService(nodes)
	if err != nil {
		return err
	}

	// Create and register nodes API
	nodesapi.New(nodes).Register(a.Router())

	return nil
}
