package http

import (
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

	alRoute := regRoute.PathPrefix("/aliases").Subrouter()
	alRoute.HandleFunc("", h.listAliases).Methods(http.MethodGet)
	alRoute.HandleFunc("/{alias_key}", h.createAlias).Methods(http.MethodPost)
	alRoute.HandleFunc("/{alias_key}", h.getAlias).Methods(http.MethodGet)
	alRoute.HandleFunc(, h.updateAlias).Methods()

	r.Methods(http.MethodDelete).Path("").HandlerFunc(h.deleteAlias)
	r.Methods(http.MethodDelete).Path("").HandlerFunc(h.deleteAlias)
	r.Methods(http.MethodDelete).Path("").HandlerFunc(h.deleteAlias)
	r.Methods(http.MethodPatch).Path("/{key}").HandlerFunc(h.deleteAlias)
	r.Methods(http.MethodDelete).Path("").HandlerFunc(h.deleteAlias)
}

// @Summary Creates an alias
// @Description Create an alias of a key in a dedicated alias registry
// @Tags Aliases
// @Accept json
// @Produce json
// @Param registry_name path string true "registry identifier"
// @Param alias_key path string true "alias identifier"
// @Param request body types.AliasRequest true "Create Alias Request"
// @Success 200 {object} types.AliasResponse "Alias data"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /registries/{registry_name}/aliases/{alias_key} [post]
func (h *AliasHandler) createAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exist in this subrouter
	regName := vars["registry_name"]
	key := vars["alias_key"]

	var aliasReq types.AliasRequest
	err := jsonutils.UnmarshalBody(r.Body, &aliasReq)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, errors.InvalidFormatError(err.Error()))
		return
	}

	alias := types.FormatAlias(regName, key, aliasReq.Kind, aliasReq.Value)
	respAlias, err := h.aliases.Create(r.Context(), alias.RegistryName, alias)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}

	resp := types.AliasResponse{
		Kind:  respAlias.Value.Kind,
		Value: respAlias.Value.Value,
	}
	err = infrahttp.WriteJSON(w, resp)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}
}

// @Summary Get an alias
// @Description Get an alias of a key from a dedicated alias registry
// @Tags Aliases
// @Produce json
// @Param registry_name path string true "registry identifier"
// @Param alias_key path string true "alias identifier"
// @Success 200 {object} types.AliasResponse "Alias data"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Alias not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /registries/{registry_name}/aliases/{alias_key} [get]
func (h *AliasHandler) getAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exist in this subrouter
	regName := vars["registry_name"]
	key := vars["alias_key"]

	alias, err := h.aliases.Get(r.Context(), regName, key)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}

	resp := types.AliasResponse{
		Kind:  alias.Value.Kind,
		Value: alias.Value.Value,
	}
	err = infrahttp.WriteJSON(w, resp)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}
}

// updateAlias updates an alias value.
// @Summary Update an alias
// @Description Update an alias of a key from a dedicated alias registry
// @Tags Aliases
// @Accept json
// @Produce json
// @Param registry_name path string true "registry identifier"
// @Param alias_key path string true "alias identifier"
// @Param request body types.AliasRequest true "Update Alias Request"
// @Success 200 {object} types.AliasResponse "Alias data"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Alias not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /registries/{registry_name}/aliases/{alias_key} [put]
func (h *AliasHandler) updateAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exist in this subrouter
	regName := vars["registry_name"]
	key := vars["alias_key"]

	var aliasReq types.AliasRequest
	err := jsonutils.UnmarshalBody(r.Body, &aliasReq)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, errors.InvalidFormatError(err.Error()))
		return
	}

	alias := types.FormatAlias(regName, key, aliasReq.Kind, aliasReq.Value)

	newAlias, err := h.aliases.Update(r.Context(), regName, alias)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}

	resp := types.AliasResponse{
		Kind:  newAlias.Value.Kind,
		Value: newAlias.Value.Value,
	}
	err = infrahttp.WriteJSON(w, resp)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
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

	err := h.aliases.Delete(ctx, RegistryNameFromContext(ctx), mux.Vars(r)["key"])
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
