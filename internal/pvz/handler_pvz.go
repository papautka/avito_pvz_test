package pvz

import (
	"avito_pvz_test/internal/dto/errorDto"
	"avito_pvz_test/internal/dto/payload"
	"avito_pvz_test/pkg/req"
	"fmt"
	"net/http"
)

type HandlerPvz interface {
	CreatePVZ() http.HandlerFunc
	CloseLastReceptionByPvz() http.HandlerFunc
}

type HandPvz struct {
	pvzService ServicePvz
}

func NewHandPvz(pvzService ServicePvz) HandlerPvz {
	return &HandPvz{
		pvzService: pvzService,
	}
}

func (pvzHandler *HandPvz) CreatePVZ() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[payload.PvzCreateRequest](&w, r)
		if err != nil {
			strError := "Ошибка возникла на этапе чтения из объекта Request. Тело ошибки"
			errorDto.ShowResponseError(&w, strError, err, http.StatusBadRequest)
			return
		}
		pvz, err := pvzHandler.pvzService.Register(body.Id, body.RegistrationDate, body.City)
		if err != nil {
			strError := "Ошибка возникла на этапе создания PVZ в базе данных. Тело ошибки"
			errorDto.ShowResponseError(&w, strError, err, http.StatusBadRequest)
			return
		}
		req.JsonResponse(&w, pvz)
	}
}

func (pvzHandler *HandPvz) CloseLastReceptionByPvz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := r.PathValue("pvzId")
		if res == "" {
			strError := "Ошибка: отсутствует pvzId в пути запроса"
			errorDto.ShowResponseError(&w, strError, http.StatusBadRequest)
			return
		}

		uid, err := pvzHandler.pvzService.ChangeStatusReceptionByPvzOnClose(res)

		if err != nil {
			strError := fmt.Sprintf("Ошибка при закрытии последней приемки: %v", err)
			errorDto.ShowResponseError(&w, strError, http.StatusBadRequest)
			return
		}

		if uid == nil {
			strError := "Приемка не найдена или уже закрыта"
			errorDto.ShowResponseError(&w, strError, http.StatusNotFound)
			return
		}

		req.JsonResponse(&w, uid)
	}
}
