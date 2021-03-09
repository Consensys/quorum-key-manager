package accountsapi

import (
	"net/http"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/core"
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
	// // Unmarshal request body
	// reqBody := new(CreateAccountRequest)
	// if err := json.Unmarshal(req.Body, reqBody); err != nil {
	// 	// Write error
	// 	return
	// }
	// 
	// // Generate internal type object
	// attr := models.Attributes{
	// 	// Enabled:  reqBody.Enabled,
	// 	// ExpireAt: reqBody.ExpireAt,
	// 	Tags:     reqBody.Tags,
	// }
	// 
	// // Execute account creation
	// account, err := h.backend.
	// 	StoreManager().
	// 	GetAccountsStore(req.Context(), reqBody.StoreName).
	// 	Create(req.Context(), attr)
	// if err != nil {
	// 	// Write error
	// 	return
	// }
	// 
	// // Create response body from interal type
	// respBody := &CreateSecretResponse{
	// 	Addr: account.Addr,
	// 	Metadata: &Metadata{
	// 		Version:   account.Medadata.Version,
	// 		CreatedAt: account.Medadata.CreatedAt,
	// 	},
	// }
	// 
	// // Write Response
	// rw.WriteHeader(http.StatusOK)
	// _ = json.NewEncoder(rw).Encode(respBody)
}
