package config

import (
	"strconv"
)

type ServerConfig struct {
	Host string
	Port int
}

const DefaultServerPort = 8000
const DefaultServerHost = "127.0.0.1"

func ServerConfigFromEnv(env Env) (ServerConfig, error) {
	var err error
	port := DefaultServerPort
	if envPort := env["SERVER_PORT"]; envPort != "" {
		port, err = strconv.Atoi(envPort)
	}

	if err != nil {
		return ServerConfig{}, err
	}

	host := env["SERVER_HOST"]
	if host == "" {
		host = DefaultServerHost
	}

	return ServerConfig{
		Host: host,
		Port: port,
	}, nil
}
