package pvz

import (
	"avito_pvz_test/config"
	"avito_pvz_test/internal/dto/payload"
	"avito_pvz_test/pkg/midware"
	"avito_pvz_test/pkg/req"
	"log"
	"net/http"
)

type PvzHandlerDependency struct {
	*PvzService
	*config.Config
}

type PvzHandler struct {
	*PvzService
	*config.Config
}

func (pvzHandler *PvzHandler) CreatePVZ() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[payload.PvzCreateRequest](&w, r)
		if err != nil {
			log.Println("CreatePVZ: функция HandleBody вернула nil", err)
			return
		}
		pvz, err := pvzHandler.PvzService.Register(body.Id, body.RegistrationDate, body.City)
		if err != nil {
			log.Println("CreatePVZ:", err)
			return
		}
		req.JsonResponse(&w, pvz)
	}
}

func NewPvzHandler(router *http.ServeMux, pvz *PvzHandlerDependency) *PvzHandler {
	pvzHandler := &PvzHandler{
		pvz.PvzService,
		pvz.Config,
	}

	router.Handle("POST /pvz", midware.CheckRoleByToken(pvzHandler.CreatePVZ()))
	return pvzHandler
}
