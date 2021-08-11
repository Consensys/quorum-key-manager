package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	jsonutils "github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/aliases/api/types"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	"github.com/gorilla/mux"
)

type AliasHandler struct {
	alias aliases.Alias
}

func NewAliasHandler(backend aliases.Alias) *AliasHandler {
	h := AliasHandler{
		alias: backend,
	}

	return &h
}

func (h *AliasHandler) Register(r *mux.Router) {
	registries := r.PathPrefix("/registries/{registry_name}").Subrouter()
	registries.HandleFunc("", h.deleteRegistry).Methods(http.MethodDelete)

	aliases := registries.PathPrefix("/aliases").Subrouter()
	aliases.HandleFunc("", h.createAlias).Methods(http.MethodPost)
	aliases.HandleFunc("", h.listAliases).Methods(http.MethodGet)
	aliases.HandleFunc("/{alias_key}", h.getAlias).Methods(http.MethodGet)
	aliases.HandleFunc("/{alias_key}", h.updateAlias).Methods(http.MethodPut)
	aliases.HandleFunc("/{alias_key}", h.deleteAlias).Methods(http.MethodDelete)
}

func (h *AliasHandler) deleteRegistry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exists in this subrouter
	regName, _ := vars["registry_name"]

	err := h.alias.DeleteRegistry(r.Context(), aliasent.RegistryName(regName))
	if err != nil {
		WriteHTTPErrorResponse(w, err)
		return
	}
}

func (h *AliasHandler) createAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exists in this subrouter
	regName, _ := vars["registry_name"]
	rName := types.RegistryName(regName)

	aliasReq := &types.CreateAliasRequest{}
	err := jsonutils.UnmarshalBody(r.Body, aliasReq)
	if err != nil {
		WriteHTTPErrorResponse(w, errors.InvalidFormatError(err.Error()))
		return
	}

	alias, err := h.alias.CreateAlias(r.Context(), aliasent.RegistryName(regName), types.ToEntityAlias(rName, aliasReq.Alias))
	if err != nil {
		WriteHTTPErrorResponse(w, err)
		return
	}

	resp := types.CreateAliasResponse{
		Alias: types.FromEntityAlias(*alias),
	}
	err = jsonWrite(w, resp)
	if err != nil {
		WriteHTTPErrorResponse(w, err)
		return
	}
}

func (h *AliasHandler) getAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exists in this subrouter
	regName, _ := vars["registry_name"]
	aliasKey, _ := vars["alias_key"]

	alias, err := h.alias.GetAlias(r.Context(), aliasent.RegistryName(regName), aliasent.AliasKey(aliasKey))
	if err != nil {
		WriteHTTPErrorResponse(w, err)
		return
	}

	err = jsonWrite(w, alias)
	if err != nil {
		WriteHTTPErrorResponse(w, err)
		return
	}
}

// updateAlias updates an alias value.
func (h *AliasHandler) updateAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exists in this subrouter
	regName, _ := vars["registry_name"]
	aliasKey, _ := vars["alias_key"]

	aliasReq := &types.UpdateAliasRequest{}
	err := jsonutils.UnmarshalBody(r.Body, aliasReq)
	if err != nil {
		WriteHTTPErrorResponse(w, errors.InvalidFormatError(err.Error()))
		return
	}

	alias := &aliasent.Alias{
		RegistryName: aliasent.RegistryName(regName),
		Key:          aliasent.AliasKey(aliasKey),
		Value:        aliasent.AliasValue(aliasReq.Value),
	}
	// TODO the: we have to either:
	// - use an ID PK in the alias table to be able to update the alias key while renaming the alias
	// - modify the UpdateAlias func to change the alias key (PK)
	// - delete + create of the new alias
	alias, err = h.alias.UpdateAlias(r.Context(), aliasent.RegistryName(regName), *alias)
	if err != nil {
		WriteHTTPErrorResponse(w, err)
		return
	}

	err = jsonWrite(w, types.UpdateAliasResponse{Value: types.AliasValue(alias.Value)})
	if err != nil {
		WriteHTTPErrorResponse(w, err)
		return
	}
}

func (h *AliasHandler) deleteAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exists in this subrouter
	regName, _ := vars["registry_name"]
	aliasKey, _ := vars["alias_key"]

	err := h.alias.DeleteAlias(r.Context(), aliasent.RegistryName(regName), aliasent.AliasKey(aliasKey))
	if err != nil {
		WriteHTTPErrorResponse(w, err)
		return
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNoContent)
}

func (h *AliasHandler) listAliases(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// should always exists in this subrouter
	regName, _ := vars["registry_name"]

	aliases, err := h.alias.ListAliases(r.Context(), aliasent.RegistryName(regName))
	if err != nil {
		WriteHTTPErrorResponse(w, err)
		return
	}

	err = jsonWrite(w, aliases)
	if err != nil {
		WriteHTTPErrorResponse(w, err)
		return
	}
}

func jsonWrite(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8;")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	return json.NewEncoder(w).Encode(data)
}

func WriteHTTPErrorResponse(rw http.ResponseWriter, err error) {
	var writeErr error
	switch {
	case errors.IsAlreadyExistsError(err):
		writeErr = writeErrorResponse(rw, http.StatusConflict, err)
	case errors.IsNotFoundError(err):
		writeErr = writeErrorResponse(rw, http.StatusNotFound, err)
	case errors.IsUnauthorizedError(err):
		writeErr = writeErrorResponse(rw, http.StatusUnauthorized, err)
	case errors.IsInvalidFormatError(err):
		writeErr = writeErrorResponse(rw, http.StatusBadRequest, err)
	case errors.IsInvalidParameterError(err), errors.IsEncodingError(err):
		writeErr = writeErrorResponse(rw, http.StatusUnprocessableEntity, err)
	case errors.IsNotImplementedError(err), errors.IsNotSupportedError(err):
		writeErr = writeErrorResponse(rw, http.StatusNotImplemented, err)
	default:
		writeErr = writeErrorResponse(rw, http.StatusInternalServerError, fmt.Errorf(internalErrMsg))
	}
	if writeErr != nil {
		// TODO the: use logger
		log.Printf("error writing the original error: %v: %v", writeErr, err)
		http.Error(rw, writeErr.Error(), http.StatusInternalServerError)
	}
}

func writeErrorResponse(w http.ResponseWriter, status int, err error) error {
	msg, e := json.Marshal(ErrorResponse{Message: err.Error(), Code: errors.FromError(err).GetCode()})
	if e != nil {
		return e
	}

	// the: should we move that to a middleware?
	w.Header().Set("Content-Type", "application/json; charset=UTF-8;")
	// the: should we use that in every API response?
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_, err = w.Write(msg)
	return err
}

const (
	internalErrMsg = "internal server error. Please ask an admin for help or try again later"
)

// ErrorResponse is the standard API error response.
// the: should we create a common lib? What format? Should every message have a potentially
// empty error info?
type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
	Code    string `json:"code,omitempty" example:"IR001"`
}
