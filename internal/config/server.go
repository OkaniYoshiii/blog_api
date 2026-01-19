package config

import (
	"strconv"
	"time"
)

type ServerConfig struct {
	Host string
	Port int
	ReadTimeout time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout time.Duration
	IdleTimeout time.Duration
}

func ServerConfigFromEnv(env Env) (ServerConfig, error) {
	port, err := strconv.Atoi(env["SERVER_PORT"])
	if err != nil {
		return ServerConfig{}, err
	}

	readTimeout, err := strconv.Atoi(env["SERVER_READ_TIMEOUT"])
	if err != nil {
		return ServerConfig{}, err
	}

	readHeaderTimeout, err := strconv.Atoi(env["SERVER_READ_HEADER_TIMEOUT"])
	if err != nil {
		return ServerConfig{}, err
	}

	writeTimeout, err := strconv.Atoi(env["SERVER_WRITE_TIMEOUT"])
	if err != nil {
		return ServerConfig{}, err
	}

	idleTimeout, err := strconv.Atoi(env["SERVER_IDLE_TIMEOUT"])
	if err != nil {
		return ServerConfig{}, err
	}

	return ServerConfig{
		Host: env["SERVER_HOST"],
		Port: port,
		ReadTimeout: time.Duration(readTimeout),
		ReadHeaderTimeout: time.Duration(readHeaderTimeout),
		WriteTimeout: time.Duration(writeTimeout),
		IdleTimeout: time.Duration(idleTimeout),
	}, nil
}
