package nodes

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/app"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/manager"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/nodes/api"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/nodes/manager"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager"
)

func RegisterService(a *app.App) error {
	// Load manifests service
	m := new(manager.Manager)
	err := a.Service(m)
	if err != nil {
		return err
	}

	// Load stores service
	s := new(storemanager.Manager)
	err = a.Service(s)
	if err != nil {
		return err
	}

	// Create and register nodes service
	nodes := nodemanager.New(*s, *m)
	err = a.RegisterService(nodes)
	if err != nil {
		return err
	}

	// Create and register nodes API
	api.New(nodes).Register(a.Router())

	return nil
}
