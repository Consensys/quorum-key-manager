package handlers

// @TODO Remove once make gen-swagger issue is identified
// https://github.com/ConsenSys/quorum-key-manager/actions/runs/1170918392
type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
	Code    string `json:"code,omitempty" example:"IR001"`
}
