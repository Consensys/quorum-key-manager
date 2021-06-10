package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const servOneID = "service_1"
const servTwoID = "service_2"

type response struct {
	Service1 string `json:"service_1,omitempty"`
	Service2 string `json:"service_2,omitempty"`
}

func TestSuccessfulLiveness(t *testing.T) {
	handler := NewHealthzHandler()

	handler.AddLivenessCheck(servOneID, func() error { return nil })
	handler.AddLivenessCheck(servTwoID, func() error { return nil })

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
	handler.LiveEndpoint(rw, req)
	result := rw.Result()
	defer result.Body.Close()
	res := &response{}
	err := parseResponseBody(result.Body, res)
	require.NoError(t, err)

	assert.Equal(t, rw.Result().StatusCode, 200)
	assert.Equal(t, res.Service1, "OK")
	assert.Equal(t, res.Service2, "OK")
}

func TestFailLiveness(t *testing.T) {
	err := fmt.Errorf("fail to start service")
	handler := NewHealthzHandler()
	handler.AddLivenessCheck(servOneID, func() error { return nil })
	handler.AddLivenessCheck(servTwoID, func() error { return err })

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
	handler.LiveEndpoint(rw, req)

	result := rw.Result()
	defer result.Body.Close()
	res := &response{}
	err = parseResponseBody(result.Body, res)
	require.NoError(t, err)

	assert.Equal(t, result.StatusCode, 503)
	assert.Equal(t, res.Service1, "OK")
	assert.Equal(t, res.Service2, "fail to start service")
}

func TestSuccessfulReadiness(t *testing.T) {
	handler := NewHealthzHandler()

	handler.AddReadinessCheck(servOneID, func() error { return nil })
	handler.AddReadinessCheck(servTwoID, func() error { return nil })

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
	handler.ReadyEndpoint(rw, req)
	result := rw.Result()
	defer result.Body.Close()
	res := &response{}
	err := parseResponseBody(result.Body, res)
	require.NoError(t, err)

	assert.Equal(t, rw.Result().StatusCode, 200)
	assert.Equal(t, res.Service1, "OK")
	assert.Equal(t, res.Service2, "OK")
}

func TestFailReadiness(t *testing.T) {
	err := fmt.Errorf("fail to start service")
	handler := NewHealthzHandler()
	handler.AddReadinessCheck(servOneID, func() error { return nil })
	handler.AddReadinessCheck(servTwoID, func() error { return err })

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
	handler.ReadyEndpoint(rw, req)

	result := rw.Result()
	defer result.Body.Close()
	res := &response{}
	err = parseResponseBody(result.Body, res)
	require.NoError(t, err)

	assert.Equal(t, result.StatusCode, 503)
	assert.Equal(t, res.Service1, "OK")
	assert.Equal(t, res.Service2, "fail to start service")
}

func parseResponseBody(body io.ReadCloser, res interface{}) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	return dec.Decode(res)
}
