package server

import (
	"log"
	"net/http"
	"os"

	"github.com/andriystech/lgc/api/handlers"
	"github.com/andriystech/lgc/api/middlewares"
	"github.com/andriystech/lgc/config"
	"github.com/andriystech/lgc/services"
	"github.com/gorilla/mux"
)

type HttpServer interface {
	Run()
}

type HttpServerContainer struct {
	tokenService     services.TokenService
	userService      services.UserService
	webSocketService services.WebSocketService
	config           *config.ServerConfig
}

func NewHttpServer(ts services.TokenService, us services.UserService, ws services.WebSocketService, cg *config.ServerConfig) HttpServer {
	return &HttpServerContainer{tokenService: ts, userService: us, webSocketService: ws, config: cg}
}

func (hsc *HttpServerContainer) Run() {
	router := mux.NewRouter()
	router.Use(middlewares.LogHttpCalls(os.Stdout))
	router.Use(middlewares.PanicAndRecover)
	router.HandleFunc("/user/active/count", handlers.ActiveConnectionsCountHandler(hsc.webSocketService)).Methods("GET")
	router.HandleFunc("/user/active", handlers.ActiveUsersHandler(hsc.webSocketService)).Methods("GET")
	router.HandleFunc("/user/login", handlers.LogInUserHandler(hsc.userService, hsc.tokenService)).Methods("POST")
	router.HandleFunc("/user", handlers.RegisterUserHandler(hsc.userService)).Methods("POST")
	router.HandleFunc("/_health", handlers.HealthCheck).Methods("GET")
	http.HandleFunc("/chat/ws.rtm.start", handlers.WSConnectHandler(hsc.webSocketService, hsc.tokenService))
	http.Handle("/", router)

	log.Printf("Server is listening %s port", hsc.config.Port)
	log.Fatal(http.ListenAndServe(hsc.config.Port, nil))
}
