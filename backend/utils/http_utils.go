package utils

import (
	"encoding/json"
	"net/http"
)

func SetJsonError(w http.ResponseWriter, err error, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := struct {
		Message string `json:"message"`
	}{
		Message: err.Error(),
	}

	return json.NewEncoder(w).Encode(response)
}

type HttpError struct {
	statusCode int
	error      error
}

const LoginErrorMessage = "input email or password is incorrect"

func NewHttpError(statusCode int, error error) *HttpError {
	return &HttpError{
		statusCode: statusCode,
		error:      error,
	}
}

func (e *HttpError) Error() string {
	return e.error.Error()
}

func (e *HttpError) GetStatusCode() int {
	return e.statusCode
}
