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
}

func NewRegistryHandler(registries aliases.Registries) *RegistryHandler {
	return &RegistryHandler{registries: registries}
}

func (h *RegistryHandler) Register(router *mux.Router) {
	registryRouter := router.PathPrefix("/registries").Subrouter()

	registryRouter.Methods(http.MethodPost).Path("/{registryName}").HandlerFunc(h.create)
	registryRouter.Methods(http.MethodGet).Path("/{registryName}").HandlerFunc(h.get)
	registryRouter.Methods(http.MethodDelete).Path("/{registryName}").HandlerFunc(h.delete)
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
	registry, err := h.registries.Create(r.Context(), getRegistry(r))
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
	registry, err := h.registries.Get(r.Context(), getRegistry(r))
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
	err := h.registries.Delete(r.Context(), getRegistry(r))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func getRegistry(r *http.Request) string {
	return mux.Vars(r)["registryName"]
}
