package http

import (
	auth "github.com/consensys/quorum-key-manager/src/auth/api/http"
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	jsonutils "github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/aliases/api/types"
	infrahttp "github.com/consensys/quorum-key-manager/src/infra/http"
	"github.com/gorilla/mux"
)

type AliasHandler struct {
	aliases aliases.Aliases
}

func NewAliasHandler(aliases aliases.Aliases) *AliasHandler {
	return &AliasHandler{aliases: aliases}
}

func (h *AliasHandler) Register(r *mux.Router) {
	aliasRouter := r.PathPrefix("/registries/{registryName}/aliases").Subrouter()

	aliasRouter.Methods(http.MethodPost).Path("/{key}").HandlerFunc(h.create)
	aliasRouter.Methods(http.MethodGet).Path("/{key}").HandlerFunc(h.get)
	aliasRouter.Methods(http.MethodPatch).Path("/{key}").HandlerFunc(h.updateAlias)
	aliasRouter.Methods(http.MethodDelete).Path("").HandlerFunc(h.delete)
}

// @Summary Creates an alias
// @Description Create an alias of a key in a dedicated alias registry
// @Tags Aliases
// @Accept json
// @Produce json
// @Param registryName path string true "registry identifier"
// @Param key path string true "alias identifier"
// @Param request body types.AliasRequest true "Create Alias Request"
// @Success 200 {object} types.AliasResponse "Alias data"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /registries/{registryName}/aliases/{key} [post]
func (h *AliasHandler) create(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	aliasReq := &types.AliasRequest{}
	err := jsonutils.UnmarshalBody(r.Body, &aliasReq)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	alias, err := h.aliases.Create(ctx, getRegistry(r), getKey(r), aliasReq.Kind, aliasReq.Value, auth.UserInfoFromContext(ctx))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = infrahttp.WriteJSON(rw, types.NewAliasResponse(alias))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}
}

// @Summary Get an alias
// @Description Get an alias by key from a dedicated alias registry
// @Tags Aliases
// @Produce json
// @Param registryName path string true "registry identifier"
// @Param key path string true "alias identifier"
// @Success 200 {object} types.AliasResponse "Alias data"
// @Failure 404 {object} ErrorResponse "Alias not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /registries/{registryName}/aliases/{key} [get]
func (h *AliasHandler) get(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	alias, err := h.aliases.Get(ctx, getRegistry(r), getKey(r), auth.UserInfoFromContext(ctx))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = infrahttp.WriteJSON(rw, types.NewAliasResponse(alias))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}
}

// @Summary Update an alias
// @Description Update an alias by key from a dedicated alias registry
// @Tags Aliases
// @Accept json
// @Produce json
// @Param registryName path string true "registry identifier"
// @Param key path string true "alias identifier"
// @Param request body types.AliasRequest true "Update Alias Request"
// @Success 200 {object} types.AliasResponse "Alias data"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Alias not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /registries/{registryName}/aliases/{key} [put]
func (h *AliasHandler) updateAlias(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	aliasReq := &types.AliasRequest{}
	err := jsonutils.UnmarshalBody(r.Body, &aliasReq)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	alias, err := h.aliases.Update(ctx, getRegistry(r), getKey(r), aliasReq.Kind, aliasReq.Value, auth.UserInfoFromContext(ctx))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = infrahttp.WriteJSON(rw, types.NewAliasResponse(alias))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}
}

// deleteAlias deletes an alias value.
// @Summary Delete an alias
// @Description Delete an alias of a key from a dedicated alias registry
// @Tags Aliases
// @Param registryName path string true "registry identifier"
// @Param key path string true "alias identifier"
// @Success 204 "Deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Alias not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /registries/{registryName}/aliases/{key} [delete]
func (h *AliasHandler) delete(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := h.aliases.Delete(ctx, getRegistry(r), getKey(r), auth.UserInfoFromContext(ctx))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func getKey(r *http.Request) string {
	return mux.Vars(r)["key"]
}
