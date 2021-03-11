package jsonrpc

import (
	"regexp"
	"strings"
)

//go:generate mockgen -source=matcher.go -destination=matcher_mock_test.go -package=jsonrpc

type Matcher interface {
	Match(*Context, *RouteMatch) bool
}

type RouteMatch struct {
	Route   *Route
	Handler Handler
	Err     error
}

type VersionMatcher struct {
	version string
}

func NewVersionMatcher(version string) *VersionMatcher {
	return &VersionMatcher{
		version: version,
	}
}

func (matcher *VersionMatcher) Match(hctx *Context, _ *RouteMatch) bool {
	return hctx.Method() == matcher.version
}

type MethodMatcher struct {
	method string
}

func NewMethodMatcher(method string) *MethodMatcher {
	return &MethodMatcher{
		method: method,
	}
}

func (matcher *MethodMatcher) Match(hctx *Context, _ *RouteMatch) bool {
	return hctx.Method() == matcher.method
}

type MethodPrefixMatcher struct {
	prefix string
}

func NewMethodPrefixMatcher(prefix string) *MethodPrefixMatcher {
	return &MethodPrefixMatcher{
		prefix: prefix,
	}
}

func (matcher *MethodPrefixMatcher) Match(hctx *Context, _ *RouteMatch) bool {
	return strings.HasPrefix(hctx.Method(), matcher.prefix)
}

type MethodRegexpMatcher struct {
	regexp *regexp.Regexp
}

func NewMethodRegexpMatcher(pattern string) (*MethodRegexpMatcher, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &MethodRegexpMatcher{
		regexp: regex,
	}, nil
}

func (matcher *MethodRegexpMatcher) Match(hctx *Context, _ *RouteMatch) bool {
	return matcher.regexp.MatchString(hctx.Method())
}
