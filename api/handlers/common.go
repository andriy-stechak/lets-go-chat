package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/andriystech/lgc/errors"
)

func ParseJsonBody(r *http.Request, v interface{}) (interface{}, *errors.AppHttpError) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Unable to read body. Reason: %s", err.Error())
		appHttpError := errors.HttpBadRequest(err.Error())
		return nil, appHttpError
	}

	if err := json.Unmarshal(body, v); err != nil {
		log.Printf("Unable to parse JSON body. Reason: %s", err.Error())
		appHttpError := errors.HttpBadRequest(err.Error())
		return nil, appHttpError
	}
	return v, nil
}

func SendJsonResponse(w http.ResponseWriter, v interface{}, status int) {
	out, err := json.Marshal(v)
	if err != nil {
		log.Printf("Unable to create response. Reason: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		appHttpError := errors.HttpInternalError(err.Error())
		appHttpError.Send(w)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	fmt.Fprint(w, string(out))
}
