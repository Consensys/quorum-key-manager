package accountsapi

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/json"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/gorilla/mux"
)

type handler struct {
	h core.Backend
}

// New creates a http.Handler to be served on /accounts
func New(bckend core.Backend) http.Handler {
	h := &handler{
		h: bckend,
	}

	router := mux.NewRouter()
	router.Methods(http.MethodPost).Path("/").HandlerFunc(h.handleCreateAccount)

	return router
}

func (h *handler) handleCreateAccount(rw http.ResponseWriter, req *http.Request) {
	// Unmarshal request body
	reqBody := new(CreateAccountRequest)
	if err := json.UnmarshalBody(req.Body, reqBody); err != nil {
		// Write error
		return
	}

	// Generate internal type object
	attr := &entities.Attributes{
		// Enabled:  reqBody.Enabled,
		// ExpireAt: reqBody.ExpireAt,
		Tags: reqBody.Tags,
	}

	// Execute account creation
	store, err := h.h.
		StoreManager().
		GetAccountStore(req.Context(), reqBody.StoreName)
	if err != nil {
		// Write error
		return
	}

	_, err = store.Create(req.Context(), attr)
	if err != nil {
		// Write error
		return
	}

	// Write response
	_, _ = rw.Write([]byte("OK"))
}
