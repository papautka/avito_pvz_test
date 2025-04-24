package req

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type FilterWithPagination struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}

func NewFilterWithPagination(startDate, endDate time.Time, page, limit int) *FilterWithPagination {
	return &FilterWithPagination{
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     page,
		Offset:    limit,
	}
}

func FilterRequest(r *http.Request) (*FilterWithPagination, error) {
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	if startDateStr == "" || endDateStr == "" || pageStr == "" || limitStr == "" {
		return nil, errors.New("нужно передать все Query параметры")
	}
	startDate, err := time.Parse("2006-01-02 15:04:05.999999+00", startDateStr)
	fmt.Println("startDate", startDate)
	if err != nil {
		return nil, errors.New("query parameters startDate передан не корректно")
	}
	endDate, err := time.Parse("2006-01-02 15:04:05.999999+00", endDateStr)
	if err != nil {
		return nil, errors.New("query parameters endDate передан не корректно")
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return nil, errors.New("page не натуральное целое число")
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, errors.New("limitStr не натуральное целое число")
	}
	if page <= 0 {
		return nil, errors.New("page должно быть натуральным числом больше 0")
	}
	if limit <= 0 {
		return nil, errors.New("limitStr должно быть натуральным числом больше 0")
	}
	// offset = (номер_страницы - 1) * размер_страницы
	offset := (page - 1) * limit
	//filter.Page = 3 (третья страница)
	//filter.Limit = 10 (по 10 записей на страницу)
	//offset := (3 - 1) * 10 = 20
	// т.е. ты пропускаешь 20 записей и начинаешь с 21

	return &FilterWithPagination{
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     limit,
		Offset:    offset,
	}, nil
}
