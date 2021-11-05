package app

import (
	"log"
	"net/http"

	"github.com/andriystech/lgc/internal/app/handlers/monitoring"
	"github.com/andriystech/lgc/internal/app/handlers/user"
	"github.com/gorilla/mux"
)

const port = ":8080"

func Run() {
	router := mux.NewRouter()
	router.HandleFunc("/user/login", user.LogInUserHandler).Methods("POST")
	router.HandleFunc("/user", user.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/_health", monitoring.HealthCheck).Methods("GET")
	http.Handle("/", router)

	log.Printf("Server is listening %s port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
