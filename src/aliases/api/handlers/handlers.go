package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"sync"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	jsonutils "github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/src/aliases/api/types"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	infrahttp "github.com/consensys/quorum-key-manager/src/infra/http"
	"github.com/gorilla/mux"
)

type AliasHandler struct {
	alias aliasent.AliasBackend
}

func NewAliasHandler(backend aliasent.AliasBackend) *AliasHandler {
	h := AliasHandler{
		alias: backend,
	}

	return &h
}

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
// @Tags registry
// @Param registry_name path string true "registry identifier"
// @Param registry_key path string true "registry identifier"
// @Param request body types.DeleteRegistryRequest true "Delete Registry Request"
// @Success 204 {object} types.DeleteRegistryResponse "Registry data"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /registryes/{registry_name}/registryes/{registry_key} [delete]
func (h *AliasHandler) deleteRegistry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exists in this subrouter
	regName := vars["registry_name"]

	err := validatePathVars(regName)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, errors.InvalidFormatError(err.Error()))
		return
	}

	err = h.alias.DeleteRegistry(r.Context(), aliasent.RegistryName(regName))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Creates an alias
// @Description Create an alias of a key in a dedicated alias registry
// @Tags alias
// @Accept json
// @Produce json
// @Param registry_name path string true "registry identifier"
// @Param alias_key path string true "alias identifier"
// @Param request body types.CreateAliasRequest true "Create Alias Request"
// @Success 200 {object} types.CreateAliasResponse "Alias data"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /aliases/{registry_name}/aliases/{alias_key} [post]
func (h *AliasHandler) createAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exists in this subrouter
	regName := vars["registry_name"]
	key := vars["alias_key"]

	err := validatePathVars(regName, key)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, errors.InvalidFormatError(err.Error()))
		return
	}

	var aliasReq types.CreateAliasRequest
	err = jsonutils.UnmarshalBody(r.Body, &aliasReq)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, errors.InvalidFormatError(err.Error()))
		return
	}

	eAlias := types.FormatAlias(types.RegistryName(regName), key, aliasReq.Value)
	alias, err := h.alias.CreateAlias(r.Context(), eAlias.RegistryName, eAlias)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}

	resp := types.CreateAliasResponse{
		Value: types.AliasValue(alias.Value),
	}
	err = jsonWrite(w, resp)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}
}

// @Summary Get an alias
// @Description Get an alias of a key from a dedicated alias registry
// @Tags alias
// @Produce json
// @Param registry_name path string true "registry identifier"
// @Param alias_key path string true "alias identifier"
// @Success 200 {object} types.GetAliasResponse "Alias data"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /aliases/{registry_name}/aliases/{alias_key} [get]
func (h *AliasHandler) getAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exists in this subrouter
	regName := vars["registry_name"]
	key := vars["alias_key"]

	err := validatePathVars(regName, key)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, errors.InvalidFormatError(err.Error()))
		return
	}

	alias, err := h.alias.GetAlias(r.Context(), aliasent.RegistryName(regName), aliasent.AliasKey(key))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}

	err = jsonWrite(w, types.GetAliasResponse{
		Value: types.AliasValue(alias.Value),
	})
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}
}

// updateAlias updates an alias value.
// @Summary Update an alias
// @Description Update an alias of a key from a dedicated alias registry
// @Tags alias
// @Accept json
// @Produce json
// @Param registry_name path string true "registry identifier"
// @Param alias_key path string true "alias identifier"
// @Param request body types.UpdateAliasRequest true "Update Alias Request"
// @Success 200 {object} types.UpdateAliasResponse "Alias data"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /aliases/{registry_name}/aliases/{alias_key} [put]
func (h *AliasHandler) updateAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exists in this subrouter
	regName := vars["registry_name"]
	key := vars["alias_key"]

	err := validatePathVars(regName, key)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, errors.InvalidFormatError(err.Error()))
		return
	}

	var aliasReq types.UpdateAliasRequest
	err = jsonutils.UnmarshalBody(r.Body, &aliasReq)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, errors.InvalidFormatError(err.Error()))
		return
	}

	alias := &aliasent.Alias{
		RegistryName: aliasent.RegistryName(regName),
		Key:          aliasent.AliasKey(key),
		Value:        aliasent.AliasValue(aliasReq.Value),
	}
	// TODO the: we have to either:
	// - use an ID PK in the alias table to be able to update the alias key while renaming the alias
	// - modify the UpdateAlias func to change the alias key (PK)
	// - delete + create of the new alias
	alias, err = h.alias.UpdateAlias(r.Context(), aliasent.RegistryName(regName), *alias)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}

	err = jsonWrite(w, types.UpdateAliasResponse{
		Value: types.AliasValue(alias.Value),
	})
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}
}

// deleteAlias deletes an alias value.
// @Summary Delete an alias
// @Description Delete an alias of a key from a dedicated alias registry
// @Tags alias
// @Param registry_name path string true "registry identifier"
// @Param alias_key path string true "alias identifier"
// @Param request body types.DeleteAliasRequest true "Delete Alias Request"
// @Success 204 "Deleted successfully"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /aliases/{registry_name}/aliases/{alias_key} [delete]
func (h *AliasHandler) deleteAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exists in this subrouter
	regName := vars["registry_name"]
	key := vars["alias_key"]

	err := validatePathVars(regName, key)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, errors.InvalidFormatError(err.Error()))
		return
	}

	err = h.alias.DeleteAlias(r.Context(), aliasent.RegistryName(regName), aliasent.AliasKey(key))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Get all the aliases in a registry
// @Description Get all the aliases in a registry
// @Tags alias
// @Produce json
// @Param registry_name path string true "registry identifier"
// @Param alias_key path string true "alias identifier"
// @Success 200 {array} types.GetAliasResponse "a list of Aliases"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /aliases/{registry_name}/aliases/{alias_key} [get]
func (h *AliasHandler) listAliases(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exists in this subrouter
	regName := vars["registry_name"]

	err := validatePathVars(regName)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, errors.InvalidFormatError(err.Error()))
		return
	}

	als, err := h.alias.ListAliases(r.Context(), aliasent.RegistryName(regName))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}

	err = jsonWrite(w, types.FormatEntityAliases(als))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(w, err)
		return
	}
}

var pathVarsRegexCompileOnce sync.Once
var pathVarsRegex *regexp.Regexp

const pathVarsFormat = "^[a-zA-Z0-9-_+]+$"

func validatePathVars(pathVars ...string) error {
	var err error
	pathVarsRegexCompileOnce.Do(func() {
		pathVarsRegex, err = regexp.Compile(pathVarsFormat)
	})
	if err != nil {
		return err
	}

	for _, v := range pathVars {
		if !pathVarsRegex.MatchString(v) {
			return fmt.Errorf("`%v` in path is not in the correct format: %v", v, pathVarsFormat)
		}
	}
	return nil
}
