package errors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AppHttpError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func HttpBadRequest(message string) *AppHttpError {
	return &AppHttpError{Status: http.StatusBadRequest, Message: message}
}

func HttpUnauthorized(message string) *AppHttpError {
	return &AppHttpError{Status: http.StatusUnauthorized, Message: message}
}

func HttpNotFound(message string) *AppHttpError {
	return &AppHttpError{Status: http.StatusNotFound, Message: message}
}

func HttpConflict(message string) *AppHttpError {
	return &AppHttpError{Status: http.StatusConflict, Message: message}
}

func HttpInternalError(message string) *AppHttpError {
	return &AppHttpError{Status: http.StatusInternalServerError, Message: message}
}

func (AppHttpError *AppHttpError) Error() string {
	return AppHttpError.Message
}

func (AppHttpError *AppHttpError) Send(w http.ResponseWriter) {
	out, err := json.Marshal(AppHttpError)
	if err != nil {
		log.Printf("Unable to create response. Reason: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(AppHttpError.Status)
	fmt.Fprint(w, string(out))
}
