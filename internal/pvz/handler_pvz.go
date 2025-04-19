package pvz

import (
	"avito_pvz_test/config"
	"avito_pvz_test/internal/dto/payload"
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
		pvzHandler.PvzService.Register(body.Id, body.RegistrationDate, body.City)
	}
}
