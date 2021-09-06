package server

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
)

type HealthzHandler struct {
	http.ServeMux
	mux       sync.RWMutex
	liveness  map[string]CheckFunc
	readiness map[string]CheckFunc
}

type CheckFunc func(context.Context) error

func NewHealthzHandler() *HealthzHandler {
	h := &HealthzHandler{
		liveness:  make(map[string]CheckFunc),
		readiness: make(map[string]CheckFunc),
	}
	h.Handle("/live", http.HandlerFunc(h.LiveEndpoint))
	h.Handle("/ready", http.HandlerFunc(h.ReadyEndpoint))
	return h
}

func (s *HealthzHandler) LiveEndpoint(w http.ResponseWriter, r *http.Request) {
	s.handle(w, r, s.liveness)
}

func (s *HealthzHandler) ReadyEndpoint(w http.ResponseWriter, r *http.Request) {
	s.handle(w, r, s.readiness, s.liveness)
}

func (s *HealthzHandler) AddLivenessCheck(name string, check CheckFunc) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.liveness[name] = check
}

func (s *HealthzHandler) AddReadinessCheck(name string, check CheckFunc) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.readiness[name] = check
}

func (s *HealthzHandler) collectChecks(ctx context.Context, checks map[string]CheckFunc, resultsOut map[string]string, statusOut *int) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	for name, check := range checks {
		if _, ok := resultsOut[name]; ok {
			continue
		}
		if err := check(ctx); err != nil {
			*statusOut = http.StatusServiceUnavailable
			resultsOut[name] = err.Error()
		} else {
			resultsOut[name] = "OK"
		}
	}
}

func (s *HealthzHandler) handle(w http.ResponseWriter, r *http.Request, checks ...map[string]CheckFunc) {
	if r.Method != http.MethodGet {
		http.Error(w, "not allowed", http.StatusMethodNotAllowed)
		return
	}

	checkResults := make(map[string]string)
	status := http.StatusOK
	for _, checks := range checks {
		s.collectChecks(r.Context(), checks, checkResults, &status)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	_ = encoder.Encode(checkResults)
}
