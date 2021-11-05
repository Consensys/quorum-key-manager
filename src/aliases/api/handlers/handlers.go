package handlers

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
	alias aliases.Interactor
}

func NewAliasHandler(alias aliases.Interactor) *AliasHandler {
	h := AliasHandler{
		alias: alias,
	}

	return &h
}

// Register registers alias handlers to HTTP endpoints.
func (h *AliasHandler) Register(r *mux.Router) {
	regRoute := r.PathPrefix("/registries/{registry_name}").Subrouter()
	regRoute.HandleFunc("", h.deleteRegistry).Methods(http.MethodDelete)

	alRoute := regRoute.PathPrefix("/aliases").Subrouter()
	alRoute.HandleFunc("", h.listAliases).Methods(http.MethodGet)
	alRoute.HandleFunc("/{alias_key}", h.createAlias).Methods(http.MethodPost)
	alRoute.HandleFunc("/{alias_key}", h.getAlias).Methods(http.MethodGet)
	alRoute.HandleFunc("/{alias_key}", h.updateAlias).Methods(http.MethodPut)
	alRoute.HandleFunc("/{alias_key}", h.deleteAlias).Methods(http.MethodDelete)
}

// @Summary Delete a registry
// @Description Delete a registry and all its keys
// @Tags Registries
// @Param registry_name path string true "registry identifier"
// @Success 204 "Deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Registry not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /registries/{registry_name} [delete]
func (h *AliasHandler) deleteRegistry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exist in this subrouter
	regName := vars["registry_name"]

	err := h.alias.DeleteRegistry(r.Context(), regName)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNoContent)
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

	eAlias := types.FormatAlias(regName, key, aliasReq.AliasValue)
	alias, err := h.alias.CreateAlias(r.Context(), eAlias.RegistryName, eAlias)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}

	resp := types.AliasResponse{
		AliasValue: types.AliasValue{
			RawKind:  alias.Value.Kind,
			RawValue: alias.Value.Value,
		},
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

	alias, err := h.alias.GetAlias(r.Context(), regName, key)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}

	resp := types.AliasResponse{
		AliasValue: types.AliasValue{
			RawKind:  alias.Value.Kind,
			RawValue: alias.Value.Value,
		},
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

	alias := types.FormatAlias(regName, key, aliasReq.AliasValue)

	newAlias, err := h.alias.UpdateAlias(r.Context(), regName, alias)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}

	resp := types.AliasResponse{
		AliasValue: types.AliasValue{
			RawKind:  newAlias.Value.Kind,
			RawValue: newAlias.Value.Value,
		},
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
// @Param registry_name path string true "registry identifier"
// @Param alias_key path string true "alias identifier"
// @Success 204 "Deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Alias not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /registries/{registry_name}/aliases/{alias_key} [delete]
func (h *AliasHandler) deleteAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exist in this subrouter
	regName := vars["registry_name"]
	key := vars["alias_key"]

	err := h.alias.DeleteAlias(r.Context(), regName, key)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Get all the aliases in a registry
// @Description Get all the aliases in a registry
// @Tags Aliases
// @Produce json
// @Param registry_name path string true "registry identifier"
// @Param alias_key path string true "alias identifier"
// @Success 200 {array} types.Alias "a list of Aliases"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /registries/{registry_name}/aliases [get]
func (h *AliasHandler) listAliases(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exist in this subrouter
	regName := vars["registry_name"]

	als, err := h.alias.ListAliases(r.Context(), regName)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}

	err = infrahttp.WriteJSON(w, types.FormatEntityAliases(als))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}
}
