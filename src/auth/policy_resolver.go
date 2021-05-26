package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/armon/go-radix"
)

var errMissingAllow = fmt.Errorf("no 'Allow' statement policy (consider adding one)")

type Operation struct {
	// Tags associated to the operation
	Tags map[string]string

	// Action to be performed
	Action string

	// Path of the operation
	Path string
}

type Result struct {
	Allowed           bool
	IsRoot            bool
	UnauthorizedError error
}

type PolicyResolver interface {
	IsAuthorized(ctx context.Context, op *Operation) (ret *Result)
}

type PathPermissions struct {
	exactActions    map[string]*ActionPermissions
	prefixedActions *radix.Tree
}

func (perms *PathPermissions) IsAuthorized(op *Operation) (err error) {
	actionPerms, hasExact := perms.exactActions[op.Action]
	if hasExact {
		err = actionPerms.IsAuthorized(op)
		if err != nil {
			return
		}
	}

	actionPrefix, _, hasPrefix := perms.prefixedActions.LongestPrefix(op.Action)
	if !hasExact && !hasPrefix {
		return errMissingAllow
	}

	perms.prefixedActions.WalkPath(actionPrefix, func(_ string, v interface{}) bool {
		err = v.(*ActionPermissions).IsAuthorized(op)
		return err != nil
	})

	return
}

type ActionPermissions struct {
	Permissions []*ActionPermission
}

func (perms *ActionPermissions) IsAuthorized(op *Operation) error {
	for _, perm := range perms.Permissions {
		err := perm.IsAuthorized(op)
		if err != nil {
			return err
		}
	}

	return nil
}

type ActionPermission struct {
	Policy    *Policy
	Statement *Statement
}

func (perm *ActionPermission) IsAuthorized(_ *Operation) error {
	switch perm.Statement.Effect {
	case "Allow":
		return nil
	case "Deny":
		return fmt.Errorf("statement %q on policy %q denied access", perm.Policy.Name, perm.Statement.Name)
	default:
		panic(fmt.Sprintf("invalid effect %q", perm.Statement.Effect))
	}
}

type ACLPolicyResolver struct {
	exactPath    map[string]*PathPermissions
	prefixedPath *radix.Tree
}

func unauthorizedOpErr(op *Operation, err error) error {
	return fmt.Errorf("operation %q on resource %q unauthorized: %v", op.Action, op.Path, err)
}

func (r *ACLPolicyResolver) IsAuthorized(ctx context.Context, op *Operation) (ret *Result) {
	ret = new(Result)

	actionPerms, hasExact := r.exactPath[op.Path]
	if hasExact {
		err := actionPerms.IsAuthorized(op)
		if err != nil && err != errMissingAllow {
			ret.UnauthorizedError = unauthorizedOpErr(op, err)
			return
		}
	}

	actionPrefix, _, hasPrefix := r.prefixedPath.LongestPrefix(op.Path)
	if !hasPrefix {
		ret.UnauthorizedError = unauthorizedOpErr(op, errMissingAllow)
		return
	}

	r.prefixedPath.WalkPath(actionPrefix, func(_ string, v interface{}) bool {
		err := v.(*PathPermissions).IsAuthorized(op)
		if err != nil {
			ret.UnauthorizedError = unauthorizedOpErr(op, err)
			return true
		}
		return false
	})

	ret.Allowed = ret.UnauthorizedError == nil

	return
}

func (r *ACLPolicyResolver) insertStatement(policy *Policy, statement *Statement) {
	insertActions := func(pathPerms *PathPermissions) {
		for _, action := range statement.Actions {
			actionPerm := &ActionPermission{
				Policy:    policy,
				Statement: statement,
			}

			if strings.HasSuffix(action, "*") {
				actionPrefix := strings.TrimSuffix(action, "*")
				perms, ok := pathPerms.prefixedActions.Get(actionPrefix)
				if ok {
					perms.(*ActionPermissions).Permissions = append(perms.(*ActionPermissions).Permissions, actionPerm)
				} else {
					_, _ = pathPerms.prefixedActions.Insert(actionPrefix, &ActionPermissions{
						Permissions: []*ActionPermission{actionPerm},
					})
				}
			} else {
				actionPerms, ok := pathPerms.exactActions[action]
				if ok {
					actionPerms.Permissions = append(actionPerms.Permissions, actionPerm)
				} else {
					pathPerms.exactActions[action] = &ActionPermissions{
						Permissions: []*ActionPermission{actionPerm},
					}
				}
			}
		}
	}

	for _, path := range statement.Resource {
		if strings.HasSuffix(path, "*") {
			pathPrefix := strings.TrimSuffix(path, "*")
			perms, ok := r.prefixedPath.Get(pathPrefix)
			if !ok {
				pathPerms := &PathPermissions{
					exactActions:    make(map[string]*ActionPermissions),
					prefixedActions: radix.New(),
				}
				_, _ = r.prefixedPath.Insert(pathPrefix, pathPerms)
				insertActions(pathPerms)
			} else {
				pathPerms := perms.(*PathPermissions)
				insertActions(pathPerms)
			}
		} else {
			pathPerms, ok := r.exactPath[path]
			if !ok {
				pathPerms = &PathPermissions{
					exactActions:    make(map[string]*ActionPermissions),
					prefixedActions: radix.New(),
				}
				r.exactPath[path] = pathPerms
			}
			insertActions(pathPerms)
		}
	}
}

func (r *ACLPolicyResolver) insertPolicy(policy *Policy) {
	for _, statement := range policy.Statements {
		r.insertStatement(policy, statement)
	}
}
