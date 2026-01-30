package config

import "strconv"

type ServerConfig struct {
	Address   string
	RateLimit float64
}

const DefaultServerAddress = "127.0.0.1:8000"

func ServerConfigFromEnv(env Env) (ServerConfig, error) {
	address := env["SERVER_ADDRESS"]
	if address == "" {
		address = DefaultServerAddress
	}

	rateLimit, err := strconv.ParseFloat(env["SERVER_RATE_LIMIT"], 64)
	if err != nil {
		return ServerConfig{}, err
	}

	return ServerConfig{
		Address:   address,
		RateLimit: rateLimit,
	}, nil
}
