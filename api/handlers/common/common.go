package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/andriystech/lgc/api/errors"
)

func ParseJsonBody(r *http.Request, v interface{}) (interface{}, *errors.AppError) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Unable to read body. Reason: %s", err.Error())
		appError := errors.BadRequest(err.Error())
		return nil, appError
	}

	if err := json.Unmarshal(body, v); err != nil {
		log.Printf("Unable to parse JSON body. Reason: %s", err.Error())
		appError := errors.BadRequest(err.Error())
		return nil, appError
	}
	return v, nil
}

func SendJsonResponse(w http.ResponseWriter, v interface{}, status int) {
	out, err := json.Marshal(v)
	if err != nil {
		log.Printf("Unable to create response. Reason: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		appError := errors.InternalError(err.Error())
		appError.Send(w)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	fmt.Fprint(w, string(out))
}
