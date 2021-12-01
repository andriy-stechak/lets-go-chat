package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func ParseJsonBody(r *http.Request, v interface{}) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Unable to read body. Reason: %s", err.Error())
		return nil, err
	}

	if err := json.Unmarshal(body, v); err != nil {
		log.Printf("Unable to parse JSON body. Reason: %s", err.Error())
		return nil, err
	}
	return v, nil
}

func SendErrorJsonResponse(w http.ResponseWriter, status int, message string) {
	sendJsonResponse(w, HttpErrorResponse{
		Status:  status,
		Message: message,
	}, status)
}

func sendJsonResponse(w http.ResponseWriter, v interface{}, status int) {
	out, err := json.Marshal(v)
	if err != nil {
		log.Printf("Unable to create response. Reason: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	fmt.Fprint(w, string(out))
}
