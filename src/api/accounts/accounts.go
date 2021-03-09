package accountsapi

import (
	"net/http"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/json"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

// New creates a http.Handler to be served on /accounts
func New(_ core.Backend) http.Handler {
	// TODO: to be implemented
	return nil
}

type Handler struct {
	backend core.Backend
}

type Metadata struct {
	Version int `json:"version"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
	PurgeAt   time.Time `json:"purgeAt"`
}

type CreateAccountRequest struct {
	StoreName string            `json:"name"`
	Enabled   bool              `json:"enabled"`
	ExpireAt  time.Time         `json:"expireAt"`
	Tags      map[string]string `json:"tags"`
}

type CreateAccountResponse struct {
	Addr string `json:"address"`

	Metadata *Metadata
}

func (h *Handler) ServeHTTPCreateAccount(rw http.ResponseWriter, req *http.Request) {
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
	store, err := h.backend.
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
}
