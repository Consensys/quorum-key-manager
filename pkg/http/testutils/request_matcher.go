package testutils

import (
	"io"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type requestMatcher struct {
	t *testing.T

	urlPath string
	body    []byte
}

func RequestMatcher(t *testing.T, urlPath string, Body []byte) gomock.Matcher {
	return &requestMatcher{
		t:       t,
		urlPath: urlPath,
		body:    Body,
	}
}

func (m requestMatcher) Matches(x interface{}) bool {
	req, ok := x.(*http.Request)
	if !ok {
		return false
	}

	b := make([]byte, req.ContentLength-1)
	_, _ = io.ReadFull(req.Body, b)

	urlMatch := assert.Equal(m.t, m.urlPath, req.URL.Path, "URL path should match")
	bodyMatch := assert.Equal(m.t, b, m.body, "Body should match")

	return urlMatch && bodyMatch
}

func (m requestMatcher) String() string {
	return ""
}
