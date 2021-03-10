package errors

type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
	Code    uint64 `json:"code,omitempty" example:"24000"`
}
