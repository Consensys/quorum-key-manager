package handlers

import (
	"net/http"
	"strconv"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	http2 "github.com/consensys/quorum-key-manager/src/infra/http"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/gorilla/mux"
)

type StoresHandler struct {
	secrets *SecretsHandler
	keys    *KeysHandler
	eth     *EthHandler
	utils   *UtilsHandler
}

// NewStoresHandler creates a http.Handler to be served on /stores
func NewStoresHandler(s stores.Manager) *StoresHandler {
	return &StoresHandler{
		secrets: NewSecretsHandler(s.Stores()),
		keys:    NewKeysHandler(s.Stores()),
		eth:     NewEthHandler(s.Stores()),
		utils:   NewUtilsHandler(s.Utilities()),
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

	// Register ethereum handler on /utilities
	utilsSubrouter := storeSubrouter.PathPrefix("/utilities").Subrouter()
	h.eth.Register(utilsSubrouter)
}

func storeSelector(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(WithStoreName(r.Context(), mux.Vars(r)["storeName"])))
	})
}

func getLimitOffset(request *http.Request) (rLimit, rOffset uint64, err error) {
	limit := request.URL.Query().Get("limit")
	page := request.URL.Query().Get("page")
	if limit == "" {
		limit = http2.DefaultPageSize
	}

	rLimit, err = strconv.ParseUint(limit, 10, 64)
	if err != nil {
		return 0, 0, errors.InvalidFormatError("invalid limit value")
	}

	iPage := uint64(0)
	rOffset = 0
	if page != "" {
		iPage, err = strconv.ParseUint(page, 10, 64)
		if err != nil {
			return 0, 0, errors.InvalidFormatError("invalid page value")
		}

		rOffset = iPage * rLimit
	}

	return rLimit, rOffset, nil
}
