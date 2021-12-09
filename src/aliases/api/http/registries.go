package http

import (
	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/aliases/api/types"
	infrahttp "github.com/consensys/quorum-key-manager/src/infra/http"
	"github.com/gorilla/mux"
	"net/http"
)

type RegistryHandler struct {
	registries aliases.Registries
	aliases    *AliasHandler
}

func NewRegistryHandler(registries aliases.Registries, aliasesHandler *AliasHandler) *RegistryHandler {
	return &RegistryHandler{registries: registries, aliases: aliasesHandler}
}

func (h *RegistryHandler) Register(router *mux.Router) {
	registryRouter := router.PathPrefix("/registries/{registryName}").Subrouter()
	registryRouter.Use(registrySelector)

	registryRouter.Methods(http.MethodPost).Path("").HandlerFunc(h.create)
	registryRouter.Methods(http.MethodGet).Path("").HandlerFunc(h.get)
	registryRouter.Methods(http.MethodDelete).Path("").HandlerFunc(h.delete)

	// Register aliases routes on /registries/{registryName}
	h.aliases.Register(registryRouter)
}

func registrySelector(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(WithRegistryName(r.Context(), mux.Vars(r)["registryName"])))
	})
}

// @Summary Creates an alias registry
// @Description Creates an alias registry
// @Tags Registries
// @Accept json
// @Produce json
// @Param registryName path string true "registry identifier"
// @Success 200 {object} types.RegistryResponse "Registry data"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /registries/{registryName} [post]
func (h *RegistryHandler) create(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	registry, err := h.registries.Create(ctx, RegistryNameFromContext(ctx))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = infrahttp.WriteJSON(rw, types.NewRegistryResponse(registry))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}
}

// @Summary Gets an alias registry
// @Description Gets an alias registry
// @Tags Registries
// @Produce json
// @Param registryName path string true "registry identifier"
// @Success 200 {array} types.Alias "a list of Aliases"
// @Failure 404 {object} ErrorResponse "Registry not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /registries/{registryName} [get]
func (h *RegistryHandler) get(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	registry, err := h.registries.Get(r.Context(), RegistryNameFromContext(ctx))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = infrahttp.WriteJSON(rw, types.NewRegistryResponse(registry))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}
}

// @Summary Deletes a registry
// @Description Deletes a registry and all its aliases
// @Tags Registries
// @Param registryName path string true "registry identifier"
// @Success 204 "Deleted successfully"
// @Failure 404 {object} ErrorResponse "Registry not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /registries/{registryName} [delete]
func (h *RegistryHandler) delete(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := h.registries.Delete(ctx, RegistryNameFromContext(ctx))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
