package app

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/utils/api/http"
	"github.com/consensys/quorum-key-manager/src/utils/service/utils"
)

func RegisterService(a *app.App, logger log.Logger) *utils.Utilities {
	// Business layer
	utilsService := utils.New(logger)

	// Service layer
	router := a.Router()
	http.NewUtilsHandler(utilsService).Register(router)

	return utilsService
}
