package errors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AppError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func BadRequest(message string) *AppError {
	return &AppError{Status: http.StatusBadRequest, Message: message}
}

func Unauthorized(message string) *AppError {
	return &AppError{Status: http.StatusUnauthorized, Message: message}
}

func NotFound(message string) *AppError {
	return &AppError{Status: http.StatusNotFound, Message: message}
}

func Conflict(message string) *AppError {
	return &AppError{Status: http.StatusConflict, Message: message}
}

func InternalError(message string) *AppError {
	return &AppError{Status: http.StatusInternalServerError, Message: message}
}

func (appError *AppError) Error() string {
	return appError.Message
}

func (appError *AppError) Send(w http.ResponseWriter) {
	out, err := json.Marshal(appError)
	if err != nil {
		log.Printf("Unable to create response. Reason: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(appError.Status)
	fmt.Fprint(w, string(out))
}
