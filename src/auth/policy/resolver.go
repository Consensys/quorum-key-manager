package policy

import (
	"fmt"
	"strings"

	"github.com/armon/go-radix"
	"github.com/consensys/quorum-key-manager/src/auth/types"
)

var (
	wildCardSuffix = "*"
)

type Operation struct {
	Action       string
	ResourcePath string
}

type Result interface {
	Allowed() bool
	Error() error
}

// Resolver is responsible to control whether an operation is authorized or not
// depending on the set of policies attached to the resolver
type Resolver interface {
	IsAuthorized(op *Operation) Result
}

type statement struct {
	types.Statement
	policy *types.Policy
}

type result struct {
	op *Operation

	allow *statement
	deny  *statement
}

func (res *result) Allowed() bool {
	return res.deny == nil && res.allow != nil
}

func (res *result) Error() error {
	if res.deny != nil {
		return fmt.Errorf("action %q on resource %q denied by policy %q statement %q", res.op.Action, res.op.ResourcePath, res.deny.policy.Name, res.deny.Name)
	}

	if res.allow == nil {
		return fmt.Errorf("action %q on resource %q not allowed", res.op.Action, res.op.ResourcePath)
	}

	return nil
}

// actionsResolver is associated to every statement path
// it aggregates all statement actions associated to the path into a radix tree
type actionsResolver struct {
	// exactActions maps an exact action name to an effect resolver
	exactActions map[string]*effectResolver

	// prefixedActions maps prefixed action to an effect resolver
	prefixedActions *radix.Tree
}

func newActionsResolver() *actionsResolver {
	return &actionsResolver{
		exactActions:    make(map[string]*effectResolver),
		prefixedActions: radix.New(),
	}
}

func (rslvr *actionsResolver) addStatement(sttmnt *statement) error {
	// loop over all actions in the statement
	for _, action := range sttmnt.Actions {
		if strings.HasSuffix(action, wildCardSuffix) {
			// action is a prefix pattern (e.g. Sign*)
			actionPrefix := strings.TrimSuffix(action, wildCardSuffix)

			// lookup for effect already associtated to this prefix
			effect, ok := rslvr.prefixedActions.Get(actionPrefix)
			if ok {
				// we already had an effect for prefix so we acculmulate the statement
				err := effect.(*effectResolver).addStatement(sttmnt)
				if err != nil {
					return err
				}
			} else {
				// we create the effect and insert it into the radix tree
				newEffect := new(effectResolver)
				err := newEffect.addStatement(sttmnt)
				if err != nil {
					return err
				}
				_, _ = rslvr.prefixedActions.Insert(actionPrefix, newEffect)
			}
		} else {
			// action is an exact name (e.g. SignTransaction)

			// lookup for effect already associated to this action
			effect, ok := rslvr.exactActions[action]
			if ok {
				// we already had an effect for prefix so we acculmulate the statement
				err := effect.addStatement(sttmnt)
				if err != nil {
					return err
				}
			} else {
				// we create the effect and insert it into the map of exact actions
				newEffect := new(effectResolver)
				err := newEffect.addStatement(sttmnt)
				if err != nil {
					return err
				}
				rslvr.exactActions[action] = newEffect
			}
		}
	}

	return nil
}

func (rslvr *actionsResolver) isAuthorized(op *Operation) (res *result) {
	res = &result{
		op: op,
	}

	// do we have an exact match for the action?
	actionEffects, hasExact := rslvr.exactActions[op.Action]
	if hasExact {
		// if we have an exact match and it resolves to Deny we return Deny
		res = actionEffects.isAuthorized(op)
		if res.deny != nil {
			return res
		}
	}

	// do we have a prefix matching the action?
	actionPrefix, _, hasPrefix := rslvr.prefixedActions.LongestPrefix(op.Action)
	if !hasPrefix {
		// no prefix so we return
		return res
	}

	// we have a prefix
	// we walk all matching prefix starting from root up to the longest matching prefix
	// if we meet a single Deny we stop walking and return a Deny
	rslvr.prefixedActions.WalkPath(actionPrefix, func(_ string, v interface{}) bool {
		r := v.(*effectResolver).isAuthorized(op)
		if r.deny != nil {
			// we met a Deny
			res = r
			return true
		}

		if res.allow == nil {
			res = r
		}

		return false
	})

	return
}

// effectResolver is associated to each pair path/action (possibly prefixes)
type effectResolver struct {
	allow *statement
	deny  *statement
}

func (rslvr *effectResolver) isAuthorized(op *Operation) (res *result) {
	return &result{
		op:    op,
		deny:  rslvr.deny,
		allow: rslvr.allow,
	}
}

// addStatement cumulates the effect of multiple statement
// if a singe statement is Deny, resolver will always resolve to Deny
// if multiple Deny are passed only the first one is kept
func (rslvr *effectResolver) addStatement(sttmnt *statement) error {
	if rslvr.deny != nil {
		return nil
	}

	switch sttmnt.Effect {
	case "Allow":
		rslvr.allow = sttmnt
		return nil
	case "Deny":
		rslvr.deny = sttmnt
		rslvr.allow = nil
		return nil
	default:
		return fmt.Errorf("invalid effect %q", sttmnt.Effect)
	}
}

// RadixResolver allows to perform authorization checks

// It is built from a set of Policy by aggregating all statements into some
// radix trees structure allowing optimized authorization checks
type RadixResolver struct {
	// exactPath maps an exact path to an action resolver
	exactPath map[string]*actionsResolver

	// prefixedPath maps prefixed path to an action resolver
	prefixedPath *radix.Tree
}

// NewRadixResolver creates a new RadixResolver
func NewRadixResolver(policies ...*types.Policy) (*RadixResolver, error) {
	rslvr := &RadixResolver{
		exactPath:    make(map[string]*actionsResolver),
		prefixedPath: radix.New(),
	}

	for _, policy := range policies {
		err := rslvr.insertPolicy(policy)
		if err != nil {
			return nil, err
		}
	}

	return rslvr, nil
}

// IsAuthorized check whether operation is authorized

// An operation is authorized if
// - no statement Deny the operation
// - at least one statement Allow the operation
func (r *RadixResolver) IsAuthorized(op *Operation) Result {
	return r.isAuthorized(op)
}

func (r *RadixResolver) isAuthorized(op *Operation) (res *result) {
	res = &result{
		op: op,
	}

	// do we have an exact match for the path?
	actionRslvr, hasExact := r.exactPath[op.ResourcePath]
	if hasExact {
		// if we have an exact match and it resolves to Deny we return Deny
		r := actionRslvr.isAuthorized(op)
		if r.deny != nil {
			return r
		}
		res = r
	}

	// do we have a prefix matching the action?
	actionPrefix, _, hasPrefix := r.prefixedPath.LongestPrefix(op.ResourcePath)
	if !hasPrefix {
		// no prefix so we return
		return
	}

	// we have a prefix
	// we walk all matching prefix starting from root up to the longest matching prefix
	// if we meet a single Deny we stop walking and return a Deny
	r.prefixedPath.WalkPath(actionPrefix, func(_ string, v interface{}) bool {
		r := v.(*actionsResolver).isAuthorized(op)
		if r.deny != nil {
			res = r
			return true
		}

		if res.allow == nil {
			res = r
		}

		return false
	})

	return
}

func (r *RadixResolver) insertStatement(sttmnt *statement) error {
	// loop over all resource path in the statement
	for _, path := range sttmnt.Resource {
		if strings.HasSuffix(path, wildCardSuffix) {
			// path is a prefix pattern (e.g. /path/to/*)
			pathPrefix := strings.TrimSuffix(path, wildCardSuffix)

			// lookup for an actions resolver already associtated to this prefix
			actionsRslvr, ok := r.prefixedPath.Get(pathPrefix)
			if ok {
				// we have an existing actions resolver so we accumulate the statement to it
				err := actionsRslvr.(*actionsResolver).addStatement(sttmnt)
				if err != nil {
					return err
				}
			} else {
				// we do not have an existing actions resolver
				// we create it and accumulate statement to it
				newActions := newActionsResolver()
				err := newActions.addStatement(sttmnt)
				if err != nil {
					return err
				}

				// insert action resolver to the associated prefix
				_, _ = r.prefixedPath.Insert(pathPrefix, newActions)
			}
		} else {
			// path is exact (e.g. /path/to/a)

			// lookup for an actions resolver already associtated to this path
			actionsRslvr, ok := r.exactPath[path]
			if ok {
				// we have an existing actions resolver so we accumulate the statement to it
				err := actionsRslvr.addStatement(sttmnt)
				if err != nil {
					return err
				}
			} else {
				// we do not have an existing actions resolver
				// we create it and accumulate statement to it
				newActions := newActionsResolver()
				err := newActions.addStatement(sttmnt)
				if err != nil {
					return err
				}

				// store actions resolver
				r.exactPath[path] = newActions
			}
		}
	}

	return nil
}

func (r *RadixResolver) insertPolicy(policy *types.Policy) error {
	for _, sttmnt := range policy.Statements {
		err := r.insertStatement(&statement{*sttmnt, policy})
		if err != nil {
			return err
		}
	}
	return nil
}
