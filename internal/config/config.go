package config

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
	CORS     CORSConfig
}

func FromEnv(env Env) (Config, error) {
	dbConfig, err := DatabaseConfigFromEnv(env)
	if err != nil {
		return Config{}, err
	}

	jwtConfig, err := JWTConfigFromEnv(env)
	if err != nil {
		return Config{}, err
	}

	serverConfig, err := ServerConfigFromEnv(env)
	if err != nil {
		return Config{}, err
	}

	corsConfig, err := CORSConfigFromEnv(env)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Database: dbConfig,
		JWT:      jwtConfig,
		Server:   serverConfig,
		CORS:     corsConfig,
	}, nil
}
