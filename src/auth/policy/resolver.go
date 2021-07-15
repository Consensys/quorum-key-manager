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
	Action string
	Path   string
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
		return fmt.Errorf("action %q on resource %q denied by policy %q statement %q", res.op.Action, res.op.Path, res.deny.policy.Name, res.deny.Name)
	}

	if res.allow == nil {
		return fmt.Errorf("action %q on resource %q not allowed", res.op.Action, res.op.Path)
	}

	return nil
}

type actionResolver struct {
	exactActions    map[string]*effectsResolver
	prefixedActions *radix.Tree
}

func newActionResolver() *actionResolver {
	return &actionResolver{
		exactActions:    make(map[string]*effectsResolver),
		prefixedActions: radix.New(),
	}
}

func (rslvr *actionResolver) addStatement(sttmnt *statement) error {
	for _, action := range sttmnt.Actions {
		effect, err := newEffectResolver(sttmnt)
		if err != nil {
			return err
		}

		if strings.HasSuffix(action, wildCardSuffix) {
			actionPrefix := strings.TrimSuffix(action, wildCardSuffix)
			effects, ok := rslvr.prefixedActions.Get(actionPrefix)
			if ok {
				effects.(*effectsResolver).addEffect(effect)
			} else {
				newEffects := new(effectsResolver)
				newEffects.addEffect(effect)
				_, _ = rslvr.prefixedActions.Insert(actionPrefix, newEffects)
			}
		} else {
			effects, ok := rslvr.exactActions[action]
			if ok {
				effects.addEffect(effect)
			} else {
				newEffects := new(effectsResolver)
				newEffects.addEffect(effect)
				rslvr.exactActions[action] = newEffects
			}
		}
	}

	return nil
}

func (rslvr *actionResolver) isAuthorized(op *Operation) (res *result) {
	res = &result{
		op: op,
	}

	actionEffects, hasExact := rslvr.exactActions[op.Action]
	if hasExact {
		res = actionEffects.isAuthorized(op)
		if res.deny != nil {
			return res
		}
	}

	actionPrefix, _, hasPrefix := rslvr.prefixedActions.LongestPrefix(op.Action)
	if !hasPrefix {
		return res
	}

	rslvr.prefixedActions.WalkPath(actionPrefix, func(_ string, v interface{}) bool {
		r := v.(*effectsResolver).isAuthorized(op)
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

type effectsResolver struct {
	effects []*effectResolver
}

func (rslvr *effectsResolver) isAuthorized(op *Operation) (res *result) {
	for _, effect := range rslvr.effects {
		r := effect.isAuthorized(op)
		if r.deny != nil {
			res = r
			return
		}

		if res == nil || res.allow == nil {
			res = r
		}
	}

	return
}

func (rslvr *effectsResolver) addEffect(effect *effectResolver) {
	rslvr.effects = append(rslvr.effects, effect)
}

type effectResolver struct {
	statement *statement
}

func newEffectResolver(sttmnt *statement) (*effectResolver, error) {
	if sttmnt == nil {
		return nil, fmt.Errorf("nil statement")
	}

	if (sttmnt.Effect != "Allow") && (sttmnt.Effect != "Deny") {
		return nil, fmt.Errorf("invalid statement effect %q", sttmnt.Effect)
	}

	return &effectResolver{
		statement: sttmnt,
	}, nil
}

func (r *effectResolver) isAuthorized(op *Operation) *result {
	res := &result{
		op: op,
	}
	switch r.statement.Effect {
	case "Allow":
		res.allow = r.statement
	case "Deny":
		res.deny = r.statement
	default:
		// this should never happen if using newStatementResolver
		panic(fmt.Sprintf("invalid effect %q", r.statement.Effect))
	}

	return res
}

type RadixResolver struct {
	exactPath    map[string]*actionResolver
	prefixedPath *radix.Tree
}

func NewRadixResolver(policies ...*types.Policy) (*RadixResolver, error) {
	rslvr := &RadixResolver{
		exactPath:    make(map[string]*actionResolver),
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

func (r *RadixResolver) IsAuthorized(op *Operation) Result {
	return r.isAuthorized(op)
}

func (r *RadixResolver) isAuthorized(op *Operation) (res *result) {
	res = &result{
		op: op,
	}

	actionRslvr, hasExact := r.exactPath[op.Path]
	if hasExact {
		r := actionRslvr.isAuthorized(op)
		if r.deny != nil {
			return r
		}
		res = r
	}

	actionPrefix, _, hasPrefix := r.prefixedPath.LongestPrefix(op.Path)
	if !hasPrefix {
		return
	}

	r.prefixedPath.WalkPath(actionPrefix, func(_ string, v interface{}) bool {
		r := v.(*actionResolver).isAuthorized(op)
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
	for _, path := range sttmnt.Resource {
		if strings.HasSuffix(path, wildCardSuffix) {
			pathPrefix := strings.TrimSuffix(path, wildCardSuffix)
			actionsRslvr, ok := r.prefixedPath.Get(pathPrefix)

			if ok {
				err := actionsRslvr.(*actionResolver).addStatement(sttmnt)
				if err != nil {
					return err
				}
			} else {
				newActions := newActionResolver()
				err := newActions.addStatement(sttmnt)
				if err != nil {
					return err
				}
				_, _ = r.prefixedPath.Insert(pathPrefix, newActions)
			}
		} else {
			actionsRslvr, ok := r.exactPath[path]
			if ok {
				err := actionsRslvr.addStatement(sttmnt)
				if err != nil {
					return err
				}
			} else {
				newActions := newActionResolver()
				err := newActions.addStatement(sttmnt)
				if err != nil {
					return err
				}
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
