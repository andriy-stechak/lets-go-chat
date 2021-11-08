package server

import (
	"log"
	"net/http"

	"github.com/andriystech/lgc/api/handlers/monitoring"
	"github.com/andriystech/lgc/api/handlers/user"
	"github.com/andriystech/lgc/config"
	"github.com/gorilla/mux"
)

func Run() {
	serverConfig := config.GetServerConfig()
	router := mux.NewRouter()
	router.HandleFunc("/user/login", user.LogInUserHandler).Methods("POST")
	router.HandleFunc("/user", user.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/_health", monitoring.HealthCheck).Methods("GET")
	http.Handle("/", router)

	log.Printf("Server is listening %s port", serverConfig.Port)
	log.Fatal(http.ListenAndServe(serverConfig.Port, nil))
}
