package secrets

import (
	"encoding/json"
	jsonutils "github.com/ConsenSysQuorum/quorum-key-manager/pkg/json"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/errors"
	storemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager"
	"github.com/gorilla/mux"
	"net/http"
)

const secretStoreNameHeader = "X-SECRET-STORE"

type Handler struct {
	storeManager storemanager.Manager
}

func New(storeManager storemanager.Manager) *Handler {
	return &Handler{
		storeManager: storeManager,
	}
}

func (c *Handler) Append(router *mux.Router) {
	router.Methods(http.MethodPost).Path("/secrets").HandlerFunc(c.create)
}

func (c *Handler) create(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	secretStoreName := request.Header.Get(secretStoreNameHeader)

	req := &CreateSecretRequest{}
	err := jsonutils.UnmarshalBody(request.Body, req)
	if err != nil {
		errors.WriteErrorResponse(rw, http.StatusBadRequest, err)
		return
	}

	secretStore, err := c.storeManager.GetSecretStore(ctx, secretStoreName)
	if err != nil {
		errors.WriteHTTPErrorResponse(rw, err)
		return
	}

	secret, err := secretStore.Set(ctx, req.ID, req.Value, req.Tags)
	if err != nil {
		errors.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(secret)
}
