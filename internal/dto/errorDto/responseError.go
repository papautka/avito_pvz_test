package errorDto

type ResponseError struct {
	Message string `json:"message"`
}

func NewResponseError(message string, err error) *ResponseError {
	return &ResponseError{message + ": " + err.Error()}
}
