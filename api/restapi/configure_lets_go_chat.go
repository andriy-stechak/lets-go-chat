// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	"github.com/andriystech/lgc/api/middlewares"
	"github.com/andriystech/lgc/api/restapi/operations"
	"github.com/andriystech/lgc/api/restapi/operations/chat"
	"github.com/andriystech/lgc/api/restapi/operations/user"
	"github.com/andriystech/lgc/config"
	"github.com/andriystech/lgc/facilities/mongo"
)

//go:generate swagger generate server --target ../../api --name LetsGoChat --spec ../../swagger.yml --principal interface{}

func configureFlags(api *operations.LetsGoChatAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.LetsGoChatAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	serverConfig := config.GetServerConfig()
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(serverConfig.DbConnectionTimeoutInSeconds),
	)
	db, err := mongo.NewClient(serverConfig)
	if err != nil {
		panic(err)
	}
	db.Connect(ctx)

	handlers := InitializeHandlers(db)

	api.UserCreateUserHandler = user.CreateUserHandlerFunc(handlers.RegisterUser)
	api.ChatGetActiveUsersHandler = chat.GetActiveUsersHandlerFunc(handlers.GetActiveUsers)
	api.ChatGetActiveUsersCountHandler = chat.GetActiveUsersCountHandlerFunc(handlers.GetActiveUsersCount)
	api.UserLoginUserHandler = user.LoginUserHandlerFunc(handlers.LoginUser)
	api.ChatWsRTMStartHandler = chat.WsRTMStartHandlerFunc(handlers.StartChat)

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {
		cancel()
		db.Disconnect(ctx)
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return middlewares.PanicAndRecover(middlewares.LogHttpCalls(os.Stdout)(handler))
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
