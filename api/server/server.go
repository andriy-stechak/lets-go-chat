package server

import (
	"log"
	"net/http"
	"os"

	"github.com/andriystech/lgc/api/handlers"
	"github.com/andriystech/lgc/api/middlewares"
	"github.com/andriystech/lgc/config"
	"github.com/andriystech/lgc/db/repositories"
	"github.com/andriystech/lgc/facilities/mongo"
	"github.com/andriystech/lgc/facilities/ws"
	"github.com/andriystech/lgc/services"
	"github.com/gorilla/mux"
)

func Run(db mongo.ClientHelper) {
	serverConfig := config.GetServerConfig()
	wsUpgrader := ws.NewUpgrader(serverConfig.WsReadBuffer, serverConfig.WsWriteBuffer)
	router := mux.NewRouter()
	router.Use(middlewares.LogHttpCalls(os.Stdout))
	router.Use(middlewares.PanicAndRecover)
	usersCollection := db.Database(serverConfig.DbName).Collection("users")
	usersRepository := repositories.NewUsersRepository(usersCollection)
	messagesCollection := db.Database(serverConfig.DbName).Collection("messages")
	messagesRepository := repositories.NewMessagesRepository(messagesCollection)
	tokensRepository := repositories.NewTokensRepository(serverConfig.TokenTTLInSeconds)
	connectionsRepository := repositories.NewConnectionsRepository()
	usersService := services.NewUserService(usersRepository)
	tokensService := services.NewTokenService(tokensRepository)
	wsService := services.NewWebSocketService(connectionsRepository, messagesRepository, usersRepository, wsUpgrader)
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
