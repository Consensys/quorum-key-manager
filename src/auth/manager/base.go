package manager

import (
	"context"
	"fmt"
	"sync"

	"github.com/consensys/quorum-key-manager/pkg/log"
	"github.com/consensys/quorum-key-manager/src/auth/types"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
)

var authKinds = []manifest.Kind{
	GroupKind,
	PolicyKind,
}

// BaseManager allow to manage Policies and Groups
type BaseManager struct {
	manifests manifestsmanager.Manager

	mux      sync.RWMutex
	policies map[string]*types.Policy
	groups   map[string]*types.Group

	sub    manifestsmanager.Subscription
	mnfsts chan []manifestsmanager.Message
}

func New(manifests manifestsmanager.Manager) *BaseManager {
	return &BaseManager{
		manifests: manifests,
		policies:  make(map[string]*types.Policy),
		groups:    make(map[string]*types.Group),
		mnfsts:    make(chan []manifestsmanager.Message),
	}
}

func (mngr *BaseManager) Start(ctx context.Context) error {
	mngr.mux.Lock()
	defer mngr.mux.Unlock()

	// Subscribe to manifest of Kind node
	sub, err := mngr.manifests.Subscribe(authKinds, mngr.mnfsts)
	if err != nil {
		return err
	}
	mngr.sub = sub

	// Start loading manifest
	go mngr.loadAll(ctx)

	return nil
}

func (mngr *BaseManager) Stop(context.Context) error {
	mngr.mux.Lock()
	defer mngr.mux.Unlock()

	if mngr.sub != nil {
		_ = mngr.sub.Unsubscribe()
	}
	close(mngr.mnfsts)
	return nil
}

func (mngr *BaseManager) Error() error {
	return nil
}

func (mngr *BaseManager) Close() error {
	return nil
}

func (mngr *BaseManager) policy(name string) (*types.Policy, error) {
	if policy, ok := mngr.policies[name]; ok {
		return policy, nil
	}

	return nil, fmt.Errorf("policy %q not found", name)
}

func (mngr *BaseManager) Policy(ctx context.Context, name string) (*types.Policy, error) {
	return mngr.policy(name)
}

func (mngr *BaseManager) Policies(context.Context) ([]string, error) {
	policies := make([]string, 0, len(mngr.policies))
	for policy := range mngr.policies {
		policies = append(policies, policy)
	}
	return policies, nil
}

func (mngr *BaseManager) group(name string) (*types.Group, error) {
	if group, ok := mngr.groups[name]; ok {
		return group, nil
	}

	return nil, fmt.Errorf("group %q not found", name)
}

func (mngr *BaseManager) Group(ctx context.Context, name string) (*types.Group, error) {
	return mngr.group(name)
}

func (mngr *BaseManager) Groups(context.Context) ([]string, error) {
	groups := make([]string, 0, len(mngr.groups))
	for group := range mngr.groups {
		groups = append(groups, group)
	}
	return groups, nil
}

func (mngr *BaseManager) loadAll(ctx context.Context) {
	for mnfsts := range mngr.mnfsts {
		for _, mnf := range mnfsts {
			_ = mngr.load(ctx, mnf.Manifest)
		}
	}
}

func (mngr *BaseManager) load(ctx context.Context, mnf *manifest.Manifest) error {
	mngr.mux.Lock()
	defer mngr.mux.Unlock()

	logger := log.FromContext(ctx).
		WithField("kind", mnf.Kind).
		WithField("name", mnf.Name)

	switch mnf.Kind {
	case GroupKind:
		err := mngr.loadGroup(mnf)
		if err != nil {
			logger.WithError(err).Errorf("could not load Group")
			return err
		}
		logger.Infof("loaded Group")
	case PolicyKind:
		err := mngr.loadPolicy(mnf)
		if err != nil {
			logger.WithError(err).Errorf("could not load Policy")
			return err
		}
		logger.Infof("loaded Policy")
	default:
		err := fmt.Errorf("invalid manifest kind %s", mnf.Kind)
		logger.WithError(err).Errorf("error starting node")
		return err
	}

	return nil
}

func (mngr *BaseManager) loadGroup(mnf *manifest.Manifest) error {
	if _, ok := mngr.groups[mnf.Name]; ok {
		return fmt.Errorf("group %q already exist", mnf.Name)
	}

	specs := new(GroupSpecs)
	if err := mnf.UnmarshalSpecs(specs); err != nil {
		return fmt.Errorf("invalid Group specs: %v", err)
	}

	mngr.groups[mnf.Name] = &types.Group{
		Name:     mnf.Name,
		Policies: specs.Policies,
	}

	return nil
}

func (mngr *BaseManager) loadPolicy(mnf *manifest.Manifest) error {
	if _, ok := mngr.policies[mnf.Name]; ok {
		return fmt.Errorf("policy %q already exist", mnf.Name)
	}

	specs := new(PolicySpecs)
	if err := mnf.UnmarshalSpecs(specs); err != nil {
		return fmt.Errorf("invalid Policy specs: %v", err)
	}

	mngr.policies[mnf.Name] = &types.Policy{
		Name:       mnf.Name,
		Statements: specs.Statements,
	}

	return nil
}
