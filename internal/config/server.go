package config

import (
	"strconv"
)

type ServerConfig struct {
	Host string
	Port int
}

const DefaultServerPort = 8000

func ServerConfigFromEnv(env Env) (ServerConfig, error) {
	var err error
	port := DefaultServerPort
	if envPort := env["SERVER_PORT"]; envPort != "" {
		port, err = strconv.Atoi(envPort)
	}

	if err != nil {
		return ServerConfig{}, err
	}

	return ServerConfig{
		Host: env["SERVER_HOST"],
		Port: port,
	}, nil
}
