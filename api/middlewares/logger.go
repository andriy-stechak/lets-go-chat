package middlewares

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/andriystech/lgc/api/handlers"
	gorilla "github.com/gorilla/handlers"
)

func LogHttpCalls(dst io.Writer) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return gorilla.LoggingHandler(dst, h)
	}
}

func PanicAndRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				handlers.SendErrorJsonResponse(w, http.StatusInternalServerError, fmt.Sprint(err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
