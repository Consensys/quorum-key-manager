package app

import (
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/utils/api/http"
	"github.com/consensys/quorum-key-manager/src/utils/service/utils"
	"github.com/gorilla/mux"
)

func RegisterService(router *mux.Router, logger log.Logger) *utils.Utilities {
	// Business layer
	utilsService := utils.New(logger)

	// Service layer
	http.NewUtilsHandler(utilsService).Register(router)

	return utilsService
}
