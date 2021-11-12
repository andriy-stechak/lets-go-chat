package server

import (
	"log"
	"net/http"

	"github.com/andriystech/lgc/api/handlers"
	"github.com/andriystech/lgc/config"
	"github.com/gorilla/mux"
)

func Run() {
	serverConfig := config.GetServerConfig()
	router := mux.NewRouter()
	router.HandleFunc("/user/login", handlers.LogInUserHandler).Methods("POST")
	router.HandleFunc("/user", handlers.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/_health", handlers.HealthCheck).Methods("GET")
	http.Handle("/", router)

	log.Printf("Server is listening %s port", serverConfig.Port)
	log.Fatal(http.ListenAndServe(serverConfig.Port, nil))
}
