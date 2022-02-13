package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/consensys/quorum-key-manager/pkg/errors"
)

const (
	internalErrMsg    = "internal server error. Please ask an admin for help or try again later"
	internalDepErrMsg = "failed dependency. Please ask an admin for help or try again later"
	DefaultPageSize   = "100"
)

func WriteHTTPErrorResponse(rw http.ResponseWriter, err error) {
	switch {
	case errors.IsAlreadyExistsError(err) || errors.IsStatusConflictError(err):
		writeErrorResponse(rw, http.StatusConflict, err)
	case errors.IsNotFoundError(err):
		writeErrorResponse(rw, http.StatusNotFound, err)
	case errors.IsUnauthorizedError(err):
		writeErrorResponse(rw, http.StatusUnauthorized, err)
	case errors.IsForbiddenError(err):
		writeErrorResponse(rw, http.StatusForbidden, err)
	case errors.IsInvalidFormatError(err):
		writeErrorResponse(rw, http.StatusBadRequest, err)
	case errors.IsTooManyRequestError(err):
		writeErrorResponse(rw, http.StatusTooManyRequests, err)
	case errors.IsInvalidParameterError(err), errors.IsEncodingError(err):
		writeErrorResponse(rw, http.StatusUnprocessableEntity, err)
	case errors.IsHashicorpVaultError(err), errors.IsAKVError(err), errors.IsDependencyFailureError(err), errors.IsAWSError(err), errors.IsPostgresError(err):
		writeErrorResponse(rw, http.StatusFailedDependency, errors.DependencyFailureError(internalDepErrMsg))
	case errors.IsNotImplementedError(err), errors.IsNotSupportedError(err):
		writeErrorResponse(rw, http.StatusNotImplemented, err)
	default:
		writeErrorResponse(rw, http.StatusInternalServerError, fmt.Errorf(internalErrMsg))
	}
}

func writeErrorResponse(rw http.ResponseWriter, status int, err error) {
	msg, e := json.Marshal(ErrorResponse{Message: err.Error(), Code: errors.FromError(err).GetCode()})
	if e != nil {
		http.Error(rw, e.Error(), status)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	_, _ = rw.Write(msg)
}

func WriteJSON(rw http.ResponseWriter, data interface{}) error {
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(data)
}

func WritePagingResponse(rw http.ResponseWriter, req *http.Request, data interface{}) error {
	var arrData []interface{}
	bData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bData, &arrData)
	if err != nil {
		return err
	}

	if arrData == nil {
		arrData = []interface{}{}
	}

	res := PageResponse{
		Data: arrData,
	}

	prevBaseURL, _ := url.Parse(req.Host)
	nextBaseURL, _ := url.Parse(req.Host)
	prevParams := req.URL.Query()
	nextParams := req.URL.Query()

	strLimit := req.URL.Query().Get("limit")
	if strLimit == "" {
		strLimit = DefaultPageSize
	}

	strPage := req.URL.Query().Get("page")

	if page, err := strconv.ParseUint(strPage, 10, 64); err == nil {
		nextParams.Set("page", fmt.Sprintf("%d", page+1))
		if page == 0 {
			prevParams = nil
		} else {
			prevParams.Set("page", fmt.Sprintf("%d", page-1))
		}
	} else {
		nextParams.Set("page", "1")
		prevParams = nil
	}

	if limit, err := strconv.ParseUint(strLimit, 10, 64); err == nil {
		if uint64(len(arrData)) < limit {
			nextParams = nil
		}
	}

	if prevParams != nil {
		prevBaseURL.RawQuery = prevParams.Encode()
		res.Paging.Previous = prevBaseURL.String()
		if req.TLS != nil {
			res.Paging.Previous = "https://" + res.Paging.Previous
		}
	}
	if nextParams != nil {
		nextBaseURL.RawQuery = nextParams.Encode()
		res.Paging.Next = nextBaseURL.String()
		if req.TLS != nil {
			res.Paging.Next = "https://" + res.Paging.Next
		}
	}

	return WriteJSON(rw, res)
}
