package policy

import (
	"fmt"
	"testing"

	"github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRadixResolver(t *testing.T) {
	type testIsAuthorized struct {
		desc string

		op *Operation

		expectedIsAllowed bool
		expectedErr       error
	}

	tests := []struct {
		desc string

		policies    []types.Policy
		expectedErr error
		tests       []*testIsAuthorized
	}{
		{
			desc: "Single policy single statement exact paths and actions ",
			policies: []types.Policy{
				types.Policy{
					Name: "TestPolicy1",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement1",
							Effect:   "Allow",
							Actions:  []string{"Action1", "Action2"},
							Resource: []string{"/path/to/a", "/path/to/b"},
						},
					},
				},
			},
			tests: []*testIsAuthorized{
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #1",
					op:                &Operation{Action: "Action1", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #2",
					op:                &Operation{Action: "Action1", ResourcePath: "/path/to/b"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #3",
					op:                &Operation{Action: "Action2", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #4",
					op:                &Operation{Action: "Action2", ResourcePath: "/path/to/b"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Not allowed Action on Allowed path",
					op:                &Operation{Action: "Action3", ResourcePath: "/path/to/a"},
					expectedIsAllowed: false,
					expectedErr:       fmt.Errorf("action \"Action3\" on resource \"/path/to/a\" not allowed"),
				},
				&testIsAuthorized{
					desc:              "Allowed Action on not allowed path",
					op:                &Operation{Action: "Action1", ResourcePath: "/path/to/c"},
					expectedIsAllowed: false,
					expectedErr:       fmt.Errorf("action \"Action1\" on resource \"/path/to/c\" not allowed"),
				},
			},
		},
		{
			desc: "Single policy multiple statements exact paths and actions",
			policies: []types.Policy{
				types.Policy{
					Name: "TestPolicy1",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement1",
							Effect:   "Allow",
							Actions:  []string{"Action1"},
							Resource: []string{"/path/to/a"},
						},
						&types.Statement{
							Name:     "TestStatement1",
							Effect:   "Allow",
							Actions:  []string{"Action2"},
							Resource: []string{"/path/to/b"},
						},
					},
				},
			},
			tests: []*testIsAuthorized{
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #1",
					op:                &Operation{Action: "Action1", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #2",
					op:                &Operation{Action: "Action2", ResourcePath: "/path/to/b"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:        "Not allowed Action/ResourcePath combination #1",
					op:          &Operation{Action: "Action2", ResourcePath: "/path/to/a"},
					expectedErr: fmt.Errorf("action \"Action2\" on resource \"/path/to/a\" not allowed"),
				},
				&testIsAuthorized{
					desc:        "Not allowed Action/ResourcePath combination #2",
					op:          &Operation{Action: "Action1", ResourcePath: "/path/to/b"},
					expectedErr: fmt.Errorf("action \"Action1\" on resource \"/path/to/b\" not allowed"),
				},
			},
		},
		{
			desc: "Multiple policy multiple statements exact paths and actions",
			policies: []types.Policy{
				types.Policy{
					Name: "TestPolicy1",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement1",
							Effect:   "Allow",
							Actions:  []string{"Action1"},
							Resource: []string{"/path/to/a"},
						},
						&types.Statement{
							Name:     "TestStatement1",
							Effect:   "Allow",
							Actions:  []string{"Action2"},
							Resource: []string{"/path/to/b"},
						},
					},
				},
				types.Policy{
					Name: "TestPolicy2",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement3",
							Effect:   "Allow",
							Actions:  []string{"Action1"},
							Resource: []string{"/path/to/b"},
						},
						&types.Statement{
							Name:     "TestStatement4",
							Effect:   "Allow",
							Actions:  []string{"Action2"},
							Resource: []string{"/path/to/a"},
						},
					},
				},
			},
			tests: []*testIsAuthorized{
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #1",
					op:                &Operation{Action: "Action1", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #2",
					op:                &Operation{Action: "Action1", ResourcePath: "/path/to/b"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #3",
					op:                &Operation{Action: "Action2", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #4",
					op:                &Operation{Action: "Action2", ResourcePath: "/path/to/b"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Not allowed Action on Allowed path",
					op:                &Operation{Action: "Action3", ResourcePath: "/path/to/a"},
					expectedIsAllowed: false,
					expectedErr:       fmt.Errorf("action \"Action3\" on resource \"/path/to/a\" not allowed"),
				},
				&testIsAuthorized{
					desc:              "Allowed Action on not allowed path",
					op:                &Operation{Action: "Action1", ResourcePath: "/path/to/c"},
					expectedIsAllowed: false,
					expectedErr:       fmt.Errorf("action \"Action1\" on resource \"/path/to/c\" not allowed"),
				},
			},
		},
		{
			desc: "Multiple statements with Deny exact paths and actions",
			policies: []types.Policy{
				types.Policy{
					Name: "TestPolicy1",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement1",
							Effect:   "Allow",
							Actions:  []string{"Action1", "Action2"},
							Resource: []string{"/path/to/a", "/path/to/b"},
						},
					},
				},
				types.Policy{
					Name: "TestPolicy2",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement3",
							Effect:   "Deny",
							Actions:  []string{"Action1"},
							Resource: []string{"/path/to/b"},
						},
						&types.Statement{
							Name:     "TestStatement4",
							Effect:   "Deny",
							Actions:  []string{"Action2"},
							Resource: []string{"/path/to/a"},
						},
					},
				},
			},
			tests: []*testIsAuthorized{
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #1",
					op:                &Operation{Action: "Action1", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #2",
					op:                &Operation{Action: "Action2", ResourcePath: "/path/to/b"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:        "Denied action/path combination #1",
					op:          &Operation{Action: "Action2", ResourcePath: "/path/to/a"},
					expectedErr: fmt.Errorf("action \"Action2\" on resource \"/path/to/a\" denied by policy \"TestPolicy2\" statement \"TestStatement4\""),
				},
				&testIsAuthorized{
					desc:        "Denied action/path combination #2",
					op:          &Operation{Action: "Action1", ResourcePath: "/path/to/b"},
					expectedErr: fmt.Errorf("action \"Action1\" on resource \"/path/to/b\" denied by policy \"TestPolicy2\" statement \"TestStatement3\""),
				},
			},
		},
		{
			desc: "Single statement exact paths and wildcard action",
			policies: []types.Policy{
				types.Policy{
					Name: "TestPolicy1",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement1",
							Effect:   "Allow",
							Actions:  []string{"ActionA*", "ActionB*"},
							Resource: []string{"/path/to/a"},
						},
					},
				},
			},
			tests: []*testIsAuthorized{
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #1",
					op:                &Operation{Action: "ActionA", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #2",
					op:                &Operation{Action: "ActionAFoo", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #3",
					op:                &Operation{Action: "ActionABar", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #1",
					op:                &Operation{Action: "ActionB", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #2",
					op:                &Operation{Action: "ActionBFoo", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #3",
					op:                &Operation{Action: "ActionBBar", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:        "Too short action",
					op:          &Operation{Action: "Action", ResourcePath: "/path/to/a"},
					expectedErr: fmt.Errorf("action \"Action\" on resource \"/path/to/a\" not allowed"),
				},
				&testIsAuthorized{
					desc:        "Invalid prefix action",
					op:          &Operation{Action: "InvalidPrefixAction", ResourcePath: "/path/to/a"},
					expectedErr: fmt.Errorf("action \"InvalidPrefixAction\" on resource \"/path/to/a\" not allowed"),
				},
				&testIsAuthorized{
					desc:        "Valid action invalid path",
					op:          &Operation{Action: "ActionA", ResourcePath: "/path/to/b"},
					expectedErr: fmt.Errorf("action \"ActionA\" on resource \"/path/to/b\" not allowed"),
				},
			},
		},
		{
			desc: "Single statement wildcard path and exact action",
			policies: []types.Policy{
				types.Policy{
					Name: "TestPolicy1",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement1",
							Effect:   "Allow",
							Actions:  []string{"Action"},
							Resource: []string{"/path/to/*"},
						},
					},
				},
			},
			tests: []*testIsAuthorized{
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #1",
					op:                &Operation{Action: "Action", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #2",
					op:                &Operation{Action: "Action", ResourcePath: "/path/to/b"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #3",
					op:                &Operation{Action: "Action", ResourcePath: "/path/to/a/b/c"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #4",
					op:                &Operation{Action: "Action", ResourcePath: "/path/to/"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:        "ResourcePath too short",
					op:          &Operation{Action: "Action", ResourcePath: "/path/to"},
					expectedErr: fmt.Errorf("action \"Action\" on resource \"/path/to\" not allowed"),
				},
				&testIsAuthorized{
					desc:        "Invalid path",
					op:          &Operation{Action: "Action", ResourcePath: "/path/too/a"},
					expectedErr: fmt.Errorf("action \"Action\" on resource \"/path/too/a\" not allowed"),
				},
				&testIsAuthorized{
					desc:        "Invalid action",
					op:          &Operation{Action: "InvalidAction", ResourcePath: "/path/to/a"},
					expectedErr: fmt.Errorf("action \"InvalidAction\" on resource \"/path/to/a\" not allowed"),
				},
			},
		},
		{
			desc: "Single statement wildcard path and wildcar action",
			policies: []types.Policy{
				types.Policy{
					Name: "TestPolicy1",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement1",
							Effect:   "Allow",
							Actions:  []string{"Action*"},
							Resource: []string{"/path/to/*"},
						},
					},
				},
			},
			tests: []*testIsAuthorized{
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #1",
					op:                &Operation{Action: "ActionA", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #2",
					op:                &Operation{Action: "ActionB", ResourcePath: "/path/to/b"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:        "ResourcePath too short",
					op:          &Operation{Action: "Action", ResourcePath: "/path/to"},
					expectedErr: fmt.Errorf("action \"Action\" on resource \"/path/to\" not allowed"),
				},
				&testIsAuthorized{
					desc:        "ResourcePath too short",
					op:          &Operation{Action: "Action", ResourcePath: "/path/to"},
					expectedErr: fmt.Errorf("action \"Action\" on resource \"/path/to\" not allowed"),
				},
				&testIsAuthorized{
					desc:        "Action too short",
					op:          &Operation{Action: "Actio", ResourcePath: "/path/to/a"},
					expectedErr: fmt.Errorf("action \"Actio\" on resource \"/path/to/a\" not allowed"),
				},
			},
		},
		{
			desc: "Multiple statement with Deny wildcard path and wildcard action",
			policies: []types.Policy{
				types.Policy{
					Name: "TestPolicy1",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement1",
							Effect:   "Allow",
							Actions:  []string{"Action*"},
							Resource: []string{"/path/to/*"},
						},
					},
				},
				types.Policy{
					Name: "TestPolicy2",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement2",
							Effect:   "Deny",
							Actions:  []string{"ActionDenied*"},
							Resource: []string{"/path/*"},
						},
						&types.Statement{
							Name:     "TestStatement3",
							Effect:   "Deny",
							Actions:  []string{"ActionADenied*"},
							Resource: []string{"/path/to/a*"},
						},
					},
				},
			},
			tests: []*testIsAuthorized{
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #1",
					op:                &Operation{Action: "ActionA", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #2",
					op:                &Operation{Action: "ActionB", ResourcePath: "/path/to/b"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:        "Action denied",
					op:          &Operation{Action: "ActionDenied", ResourcePath: "/path/to/a"},
					expectedErr: fmt.Errorf("action \"ActionDenied\" on resource \"/path/to/a\" denied by policy \"TestPolicy2\" statement \"TestStatement2\""),
				},
				&testIsAuthorized{
					desc:              "ActionAdenied on not denied on path",
					op:                &Operation{Action: "ActionADenied", ResourcePath: "/path/to/b"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:        "ActionAdenied on denied path",
					op:          &Operation{Action: "ActionADenied", ResourcePath: "/path/to/a"},
					expectedErr: fmt.Errorf("action \"ActionADenied\" on resource \"/path/to/a\" denied by policy \"TestPolicy2\" statement \"TestStatement3\""),
				},
				&testIsAuthorized{
					desc:        "ActionDenied on denied path",
					op:          &Operation{Action: "ActionDenied", ResourcePath: "/path/to"},
					expectedErr: fmt.Errorf("action \"ActionDenied\" on resource \"/path/to\" denied by policy \"TestPolicy2\" statement \"TestStatement2\""),
				},
				&testIsAuthorized{
					desc:        "ActionDenied not allowed on short path",
					op:          &Operation{Action: "ActionDenied", ResourcePath: "/path"},
					expectedErr: fmt.Errorf("action \"ActionDenied\" on resource \"/path\" not allowed"),
				},
			},
		},
		{
			desc: "Overlapping wildcards path and wildcard actions",
			policies: []types.Policy{
				types.Policy{
					Name: "TestPolicy1",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement1",
							Effect:   "Allow",
							Actions:  []string{"*"},
							Resource: []string{"/path/to/*"},
						},
						&types.Statement{
							Name:     "TestStatement2",
							Effect:   "Allow",
							Actions:  []string{"ActionAllowed"},
							Resource: []string{"*"},
						},
					},
				},
			},
			tests: []*testIsAuthorized{
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #1",
					op:                &Operation{Action: "ActionAllowed", ResourcePath: "/path/to"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #2",
					op:                &Operation{Action: "AnyAction", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:        "ResourcePath too short",
					op:          &Operation{Action: "AnyAction", ResourcePath: "/path/to"},
					expectedErr: fmt.Errorf("action \"AnyAction\" on resource \"/path/to\" not allowed"),
				},
			},
		},
		{
			desc: "Overlapping wildcards path and wildcard actions with Deny",
			policies: []types.Policy{
				types.Policy{
					Name: "TestPolicy1",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement1",
							Effect:   "Allow",
							Actions:  []string{"*"},
							Resource: []string{"/path/to/*"},
						},
						&types.Statement{
							Name:     "TestStatement2",
							Effect:   "Deny",
							Actions:  []string{"ActionDenied"},
							Resource: []string{"*"},
						},
					},
				},
			},
			tests: []*testIsAuthorized{
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #1",
					op:                &Operation{Action: "Action", ResourcePath: "/path/to/a"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:        "Allowed Action on Allowed path #2",
					op:          &Operation{Action: "ActionDenied", ResourcePath: "/path/to/a"},
					expectedErr: fmt.Errorf("action \"ActionDenied\" on resource \"/path/to/a\" denied by policy \"TestPolicy1\" statement \"TestStatement2\""),
				},
			},
		},
		{
			desc: "Same wildcard path, wildcard actions with Deny",
			policies: []types.Policy{
				types.Policy{
					Name: "TestPolicy1",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement1",
							Effect:   "Allow",
							Actions:  []string{"*"},
							Resource: []string{"*"},
						},
					},
				},
				types.Policy{
					Name: "TestPolicy2",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement2",
							Effect:   "Deny",
							Actions:  []string{"ActionDenied*"},
							Resource: []string{"*"},
						},
					},
				},
			},
			tests: []*testIsAuthorized{
				&testIsAuthorized{
					desc:              "Allowed Action on Allowed path #1",
					op:                &Operation{Action: "Action", ResourcePath: "/path"},
					expectedIsAllowed: true,
				},
				&testIsAuthorized{
					desc:        "Allowed Action on Allowed path #2",
					op:          &Operation{Action: "ActionDenied", ResourcePath: "/path"},
					expectedErr: fmt.Errorf("action \"ActionDenied\" on resource \"/path\" denied by policy \"TestPolicy2\" statement \"TestStatement2\""),
				},
			},
		},
		{
			desc: "Shorted wildcard denies first",
			policies: []types.Policy{
				types.Policy{
					Name: "TestPolicy1",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement1",
							Effect:   "Deny",
							Actions:  []string{"ABC*"},
							Resource: []string{"/a*"},
						},
					},
				},
				types.Policy{
					Name: "TestPolicy2",
					Statements: []*types.Statement{
						&types.Statement{
							Name:     "TestStatement3",
							Effect:   "Deny",
							Actions:  []string{"A*"},
							Resource: []string{"/a/b/c*"},
						},
						&types.Statement{
							Name:     "TestStatement2",
							Effect:   "Deny",
							Actions:  []string{"AB*"},
							Resource: []string{"/a/b*"},
						},
					},
				},
			},
			tests: []*testIsAuthorized{
				&testIsAuthorized{
					desc:        "ABC /a/b/c",
					op:          &Operation{Action: "ABC", ResourcePath: "/a/b/c"},
					expectedErr: fmt.Errorf("action \"ABC\" on resource \"/a/b/c\" denied by policy \"TestPolicy1\" statement \"TestStatement1\""),
				},
				&testIsAuthorized{
					desc:        "AB /a/b/c",
					op:          &Operation{Action: "AB", ResourcePath: "/a/b/c"},
					expectedErr: fmt.Errorf("action \"AB\" on resource \"/a/b/c\" denied by policy \"TestPolicy2\" statement \"TestStatement2\""),
				},
				&testIsAuthorized{
					desc:        "AB /a/b",
					op:          &Operation{Action: "AB", ResourcePath: "/a/b"},
					expectedErr: fmt.Errorf("action \"AB\" on resource \"/a/b\" denied by policy \"TestPolicy2\" statement \"TestStatement2\""),
				},
				&testIsAuthorized{
					desc:        "A",
					op:          &Operation{Action: "A", ResourcePath: "/a/b/c"},
					expectedErr: fmt.Errorf("action \"A\" on resource \"/a/b/c\" denied by policy \"TestPolicy2\" statement \"TestStatement3\""),
				},
			},
		},
	}

	for _, tts := range tests {
		t.Run(tts.desc, func(t *testing.T) {
			rslvr, err := NewRadixResolver(tts.policies...)
			if tts.expectedErr != nil {
				require.Error(t, err, "Constructing RadixResolver must fail")
				assert.Equal(t, tts.expectedErr, err, "RadixResolver construction error should be correct")
			} else {
				require.NoError(t, err, "Constructing RadixResolver must not fail")
				for _, tt := range tts.tests {
					t.Run(tt.desc, func(t *testing.T) {
						res := rslvr.IsAuthorized(tt.op)
						assert.Equal(t, tt.expectedIsAllowed, res.Allowed(), "Allowed should match")
						assert.Equal(t, tt.expectedErr, res.Error(), "Error should match")
					})
				}
			}
		})
	}
}
