package config

type ServerConfig struct {
	Address string
}

const DefaultServerAddress = "127.0.0.1:8000"

func ServerConfigFromEnv(env Env) (ServerConfig, error) {
	address := env["SERVER_ADDRESS"]
	if address == "" {
		address = DefaultServerAddress
	}

	return ServerConfig{
		Address: address,
	}, nil
}
