package errorDto

import (
	"avito_pvz_test/pkg/req"
	"fmt"
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
}

func NewResponseError(message string, err ...error) *ResponseError {
	if len(err) == 0 {
		return &ResponseError{message}
	}
	if len(err) == 1 {
		return &ResponseError{message + ": " + err[0].Error()}
	}
	return nil
}

func ShowResponseError(w *http.ResponseWriter, msg string, args ...interface{}) {
	var respErr *ResponseError
	var statusCode int

	if len(args) == 1 {
		// ShowResponseError(w, msg, statusCode)
		if code, ok := args[0].(int); ok {
			respErr = NewResponseError(msg)
			statusCode = code
		}
	} else if len(args) == 2 {
		// ShowResponse(w, msg, statusCode, err)
		if code, ok := args[0].(int); ok {
			if err, ok := args[1].(error); ok {
				respErr = NewResponseError(msg, err)
				statusCode = code
			}
		}
		fmt.Println("args[0]", args[0])
	}
	if respErr == nil {
		respErr = NewResponseError("Internal server error")
		statusCode = http.StatusInternalServerError
	}
	(*w).WriteHeader(statusCode)
	req.JsonResponse(w, respErr)
}
