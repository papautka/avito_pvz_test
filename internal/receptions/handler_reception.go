package receptions

import (
	"avito_pvz_test/internal/dto/errorDto"
	"avito_pvz_test/internal/dto/payload"
	"avito_pvz_test/pkg/req"
	"net/http"
)

type HandlerReception interface {
	CreateReception() http.HandlerFunc
}

type HandReception struct {
	serviceReception ServiceReception
}

func NewReceptionHandler(serviceReception ServiceReception) HandlerReception {
	return &HandReception{
		serviceReception: serviceReception,
	}
}

func (receptionHandler *HandReception) CreateReception() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[payload.ReceptionCreateRequest](&w, r)
		if err != nil {
			strError := "Ошибка возникла на этапе чтения из объекта Request. Тело ошибки"
			errorDto.ShowResponseError(&w, strError, err, http.StatusBadRequest)
			return
		}
		reception, err := receptionHandler.serviceReception.CreateReception(body.PvzId)
		if err != nil {
			strError := "Ошибка возникла на этапе создания reception в базе данных. Тело ошибки: " + err.Error()
			errorDto.ShowResponseError(&w, strError, http.StatusBadRequest)
			return
		}
		req.JsonResponse(&w, reception)
	}
}
