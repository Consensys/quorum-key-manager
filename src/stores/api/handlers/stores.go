package handlers

import (
	"net/http"
	"strconv"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/gorilla/mux"
)

const DefaultPageSize = "100"

type StoresHandler struct {
	stores stores.Manager

	secrets *SecretsHandler
	keys    *KeysHandler
	eth     *EthHandler
}

// NewStoresHandler creates a http.Handler to be served on /stores
func NewStoresHandler(s stores.Manager) *StoresHandler {
	return &StoresHandler{
		stores:  s,
		secrets: NewSecretsHandler(s),
		keys:    NewKeysHandler(s),
		eth:     NewAccountsHandler(s),
	}
}

func (h *StoresHandler) Register(router *mux.Router) {
	// Create subrouter for /stores
	storesSubrouter := router.PathPrefix("/stores").Subrouter()

	// Create subrouter for /stores/{storeName}
	storeSubrouter := storesSubrouter.PathPrefix("/{storeName}").Subrouter()
	storeSubrouter.Use(storeSelector)

	// Register secrets handler on /stores/{storeName}/secrets
	secretsSubrouter := storeSubrouter.PathPrefix("/secrets").Subrouter()
	h.secrets.Register(secretsSubrouter)

	// Register keys handler on /stores/{storeName}/keys
	keysSubrouter := storeSubrouter.PathPrefix("/keys").Subrouter()
	h.keys.Register(keysSubrouter)

	// Register ethereum handler on /stores/{storeName}/ethereum
	ethSubrouter := storeSubrouter.PathPrefix("/ethereum").Subrouter()
	h.eth.Register(ethSubrouter)
}

func storeSelector(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(WithStoreName(r.Context(), mux.Vars(r)["storeName"])))
	})
}

func getLimitOffset(request *http.Request) (rLimit, rOffset uint64, err error) {
	limit := request.URL.Query().Get("limit")
	page := request.URL.Query().Get("page")
	if limit == "" && page == "" {
		return 0, 0, nil
	}

	if limit == "" {
		limit = DefaultPageSize
	}

	if limit == "" {
		return 0, 0, nil
	}

	rLimit, err = strconv.ParseUint(limit, 10, 32)
	if err != nil {
		return 0, 0, errors.InvalidFormatError("invalid limit value")
	}

	iPage := uint64(0)
	rOffset = 0
	if page != "" {
		iPage, err = strconv.ParseUint(page, 10, 32)
		if err != nil {
			return 0, 0, errors.InvalidFormatError("invalid page value")
		}

		rOffset = iPage * rLimit
	}

	return rLimit, rOffset, nil
}
