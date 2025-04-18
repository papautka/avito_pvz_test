package req

import (
	"avito_pvz_test/internal/dto/errorDto"
	"encoding/json"
	"io"
	"net/http"
)

// функция которая проверяет корректность тела Request
func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		dataErr := errorDto.NewResponseError("HandleBody: Не удалось декодировать Body", err)
		(*w).WriteHeader(http.StatusBadRequest)
		JsonResponse(w, dataErr)
		return nil, err
	}
	return body, nil
}

// Функция которая выводит в формате Json ответ в Response
func JsonResponse(w *http.ResponseWriter, data any) {
	(*w).Header().Set("Content-Type", "application/json")
	json.NewEncoder(*w).Encode(data)
}

func Decode[T any](body io.ReadCloser) (*T, error) {
	// 1. Создаем структуру куда будем класть наши данные
	var payload T
	// 2. Декодируем её из JSON --> в СТРУКТУРУ
	err := json.NewDecoder(body).Decode(&payload)
	// 2.1 Если не удалось декодировать
	if err != nil {
		return nil, err
	}
	// 3. Если удалось
	return &payload, nil
}
