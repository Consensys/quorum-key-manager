package accountsapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/core"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/types"
)

// New creates a http.Handler to be served on /accounts
func New(bcknd core.Backend) http.Handler {
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
	if err := json.Unmarshal(req.Body, reqBody); err != nil {
		// Write error
		return
	}

	// Execute account creation
	attr := types.Attributes{
		Enabled:  reqBody.Enabled,
		ExpireAt: reqBody.ExpireAt,
		Tags:     reqBody.Tags,
	}
	account, err := h.backend.StoreManager().GetAccountsStore(reqBody.name).Create(req.Context(), attr)
	if err != nil {
		// Write error
		return
	}

	// Create response body
	respBody := &CreateSecretResponse{
		Addr: account.Addr,
		Metadata: &Metadata{
			Version:   account.Medadata.Version,
			CreatedAt: account.Medadata.CreatedAt,
		},
	}

	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(respBody)
}
