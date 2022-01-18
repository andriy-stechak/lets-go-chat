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
	ts services.TokenService
	us services.UserService
	ws services.WebSocketService
	cg *config.ServerConfig
}

func NewHttpServer(
	ts services.TokenService,
	us services.UserService,
	ws services.WebSocketService,
	cg *config.ServerConfig,
) HttpServer {
	return &HttpServerContainer{
		ts,
		us,
		ws,
		cg,
	}
}

func (hsc *HttpServerContainer) Run() {
	router := mux.NewRouter()
	router.Use(middlewares.LogHttpCalls(os.Stdout))
	router.Use(middlewares.PanicAndRecover)
	router.HandleFunc("/user/active/count", handlers.ActiveConnectionsCountHandler(hsc.ws)).Methods("GET")
	router.HandleFunc("/user/active", handlers.ActiveUsersHandler(hsc.ws)).Methods("GET")
	router.HandleFunc("/user/login", handlers.LogInUserHandler(hsc.us, hsc.ts)).Methods("POST")
	router.HandleFunc("/user", handlers.RegisterUserHandler(hsc.us)).Methods("POST")
	router.HandleFunc("/_health", handlers.HealthCheck).Methods("GET")
	http.HandleFunc("/chat/ws.rtm.start", handlers.WSConnectHandler(hsc.ws, hsc.ts))
	http.Handle("/", router)

	log.Printf("Server is listening %s port", hsc.cg.Port)
	log.Fatal(http.ListenAndServe(hsc.cg.Port, nil))
}
