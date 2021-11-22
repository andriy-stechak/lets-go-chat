package server

import (
	"log"
	"net/http"

	"github.com/andriystech/lgc/api/handlers"
	"github.com/andriystech/lgc/config"
	"github.com/andriystech/lgc/db/repositories"
	"github.com/andriystech/lgc/services"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func Run(db *mongo.Client) {
	serverConfig := config.GetServerConfig()
	router := mux.NewRouter()
	usersService := services.NewUserService(repositories.NewUsersRepository(db))
	router.HandleFunc("/user/login", handlers.LogInUserHandler(usersService)).Methods("POST")
	router.HandleFunc("/user", handlers.RegisterUserHandler(usersService)).Methods("POST")
	router.HandleFunc("/_health", handlers.HealthCheck).Methods("GET")
	http.Handle("/", router)

	log.Printf("Server is listening %s port", serverConfig.Port)
	log.Fatal(http.ListenAndServe(serverConfig.Port, nil))
}
