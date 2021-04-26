package jsonrpc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type routeTest struct {
	desc        string      // desc of the test
	route       *Route      // the route being tested
	msg         *RequestMsg // a request to test the route
	shouldMatch bool        // whether the request is expected to match the route at all
}

func TestVersion(t *testing.T) {
	tests := []routeTest{
		{
			desc:        "match",
			route:       new(Route).Version("2.0"),
			msg:         (&RequestMsg{}).WithVersion("2.0").WithID("abcd").WithMethod("testMethod"),
			shouldMatch: true,
		},
		{
			desc:        "not match",
			route:       new(Route).Version("2.0"),
			msg:         (&RequestMsg{}).WithVersion("3.0").WithID("abcd").WithMethod("testMethod"),
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var match RouteMatch
			if tt.shouldMatch {
				require.True(t, tt.route.Match(tt.msg, &match), "Should match")
			} else {
				require.False(t, tt.route.Match(tt.msg, &match), "Should not match")
			}
		})
	}
}

func TestMethod(t *testing.T) {
	tests := []routeTest{
		{
			desc:        "match",
			route:       new(Route).Method("knownMethod"),
			msg:         (&RequestMsg{}).WithVersion("2.0").WithID("abcd").WithMethod("knownMethod"),
			shouldMatch: true,
		},
		{
			desc:        "not match",
			route:       new(Route).Version("knownMethod"),
			msg:         (&RequestMsg{}).WithVersion("3.0").WithID("abcd").WithMethod("unknownMethod"),
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var match RouteMatch
			if tt.shouldMatch {
				require.True(t, tt.route.Match(tt.msg, &match), "Should match")
			} else {
				require.False(t, tt.route.Match(tt.msg, &match), "Should not match")
			}
		})
	}
}

func TestMethodPrefix(t *testing.T) {
	tests := []routeTest{
		{
			desc:        "match",
			route:       new(Route).MethodPrefix("known_"),
			msg:         (&RequestMsg{}).WithVersion("2.0").WithID("abcd").WithMethod("known_testMethod"),
			shouldMatch: true,
		},
		{
			desc:        "not match",
			route:       new(Route).MethodPrefix("known_"),
			msg:         (&RequestMsg{}).WithVersion("3.0").WithID("abcd").WithMethod("unknown_testMethod"),
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var match RouteMatch
			if tt.shouldMatch {
				require.True(t, tt.route.Match(tt.msg, &match), "Should match")
			} else {
				require.False(t, tt.route.Match(tt.msg, &match), "Should not match")
			}
		})
	}
}
