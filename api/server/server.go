package server

import (
	"log"
	"net/http"
	"os"

	"github.com/andriystech/lgc/api/handlers"
	"github.com/andriystech/lgc/api/middlewares"
	"github.com/andriystech/lgc/config"
	"github.com/andriystech/lgc/db/repositories"
	"github.com/andriystech/lgc/services"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func Run(db *mongo.Client) {
	serverConfig := config.GetServerConfig()
	router := mux.NewRouter()
	router.Use(middlewares.LogHttpCalls(os.Stdout))
	router.Use(middlewares.PanicAndRecover)
	usersService := services.NewUserService(repositories.NewUsersRepository(db))
	tokensService := services.NewTokenService(repositories.NewTokensRepository())
	wsService := services.NewWebSocketService(repositories.NewConnectionsRepository())
	router.HandleFunc("/user/active/count", handlers.ActiveConnectionsCountHandler(wsService)).Methods("GET")
	router.HandleFunc("/user/active", handlers.ActiveUsersHandler(wsService)).Methods("GET")
	router.HandleFunc("/user/login", handlers.LogInUserHandler(usersService, tokensService)).Methods("POST")
	router.HandleFunc("/user", handlers.RegisterUserHandler(usersService)).Methods("POST")
	router.HandleFunc("/_health", handlers.HealthCheck).Methods("GET")
	http.HandleFunc("/chat/ws.rtm.start", handlers.WSConnectHandler(wsService, tokensService))
	http.Handle("/", router)

	log.Printf("Server is listening %s port", serverConfig.Port)
	log.Fatal(http.ListenAndServe(serverConfig.Port, nil))
}
