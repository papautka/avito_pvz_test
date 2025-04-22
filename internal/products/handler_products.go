package products

import (
	"avito_pvz_test/internal/dto/errorDto"
	"avito_pvz_test/internal/dto/payload"
	"avito_pvz_test/pkg/req"
	"net/http"
)

type HandlerProduct interface {
	Create() http.HandlerFunc
}
type HandProduct struct {
	serviceProducts ServiceProduct
}

func NewHandProduct(serviceProduct ServiceProduct) HandlerProduct {
	return &HandProduct{
		serviceProducts: serviceProduct,
	}
}

func (productHandler *HandProduct) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[payload.ProductCreateRequest](&w, r)
		if err != nil {
			strError := "Ошибка возникла на этапе чтения из объекта Request. Тело ошибки" + err.Error()
			errorDto.ShowResponseError(&w, strError, http.StatusBadRequest)
			return
		}
		createdProduct, err := productHandler.serviceProducts.Create(body.PvzId, body.Type)
		if err != nil {
			strError := "Ошибка возникла на этапе создания Product в базе данных. Тело ошибки: " + err.Error()
			errorDto.ShowResponseError(&w, strError, http.StatusBadRequest)
			return
		}
		req.JsonResponse(&w, &createdProduct)
	}
}
