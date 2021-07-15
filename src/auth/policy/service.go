package policy

import (
	"context"
	"fmt"
	"sync"

	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/src/auth/types"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
)

//the: avoid globals?
var authKinds = []manifest.Kind{
	GroupKind,
	Kind,
}

// BaseManager allow to manage Policies and Groups
type BaseManager struct {
	manifests manifestsmanager.Manager

	mux      sync.RWMutex
	policies map[string]*types.Policy
	groups   map[string]*types.Group

	sub    manifestsmanager.Subscription
	mnfsts chan []manifestsmanager.Message

	logger log.Logger
}

func New(manifests manifestsmanager.Manager, logger log.Logger) *BaseManager {
	return &BaseManager{
		manifests: manifests,
		policies:  make(map[string]*types.Policy),
		groups:    make(map[string]*types.Group),
		mnfsts:    make(chan []manifestsmanager.Message),
		logger:    logger,
	}
}

func (mngr *BaseManager) Start(ctx context.Context) error {
	mngr.mux.Lock()
	defer mngr.mux.Unlock()

	// Subscribe to manifest of Kind Group and Policy
	mngr.sub = mngr.manifests.Subscribe(authKinds, mngr.mnfsts)

	// Start loading manifest
	go mngr.loadAll(ctx)

	return nil
}

func (mngr *BaseManager) Stop(context.Context) error {
	mngr.mux.Lock()
	defer mngr.mux.Unlock()

	if mngr.sub != nil {
		//the: why ignoring the error?
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

func (mngr *BaseManager) UserPolicies(ctx context.Context, info *types.UserInfo) []types.Policy {
	// Retrieve policies associated to user info
	var policies []types.Policy
	if info == nil {
		return policies
	}

	for _, groupName := range info.Groups {
		group, err := mngr.Group(ctx, groupName)
		if err != nil {
			mngr.logger.WithError(err).With("group", groupName).Debug("could not load group")
			continue
		}

		for _, policyName := range group.Policies {
			policy, err := mngr.Policy(ctx, policyName)
			if err != nil {
				mngr.logger.WithError(err).With("policy", groupName).Debug("could not load policy")
				continue
			}
			policies = append(policies, *policy)
		}
	}

	// Create resolver
	return policies
}

//the: indicate that this is not thread-safe and we need to lock the mutex before accessing this func, or lock/unlock in here.
func (mngr *BaseManager) policy(name string) (*types.Policy, error) {
	//the: https://dave.cheney.net/practical-go/presentations/gophercon-singapore-2019.html#_return_early_rather_than_nesting_deeply
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

//the: indicate that this is not thread-safe and we need to lock the mutex before accessing this func, or lock/unlock in here.
func (mngr *BaseManager) group(name string) (*types.Group, error) {
	//the: https://dave.cheney.net/practical-go/presentations/gophercon-singapore-2019.html#_return_early_rather_than_nesting_deeply
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

//the: should return an error
func (mngr *BaseManager) loadAll(ctx context.Context) {
	for mnfsts := range mngr.mnfsts {
		for _, mnf := range mnfsts {
			//the: why ignore the error?
			_ = mngr.load(ctx, mnf.Manifest)
		}
	}
}

func (mngr *BaseManager) load(_ context.Context, mnf *manifest.Manifest) error {
	//the: maybe lock in the underlying funcs (loadGroup, loadPolicy)?
	mngr.mux.Lock()
	defer mngr.mux.Unlock()

	logger := mngr.logger.With("kind", mnf.Kind).With("name", mnf.Name)

	switch mnf.Kind {
	case GroupKind:
		err := mngr.loadGroup(mnf)
		if err != nil {
			logger.WithError(err).Error("could not load Group")
			return err
		}
		logger.Info("loaded Group")
	case Kind:
		err := mngr.loadPolicy(mnf)
		if err != nil {
			logger.WithError(err).Error("could not load Policy")
			return err
		}
		logger.Info("loaded Policy")
	default:
		err := fmt.Errorf("invalid manifest kind %s", mnf.Kind)
		logger.WithError(err).Error("error starting node")
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

	specs := new(Specs)
	if err := mnf.UnmarshalSpecs(specs); err != nil {
		return fmt.Errorf("invalid Policy specs: %v", err)
	}

	mngr.policies[mnf.Name] = &types.Policy{
		Name:       mnf.Name,
		Statements: specs.Statements,
	}

	return nil
}
