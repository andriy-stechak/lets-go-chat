package config

import "os"

type ServerConfig struct {
	Port string
}

const defaultPort = ":8090"

var Server *ServerConfig

func init() {
	Server = &ServerConfig{
		Port: getPort(),
	}
}

func getPort() string {
	port, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		port = defaultPort
	}
	return port
}
