package config

type Config struct {
	Database DatabaseConfig
}

func FromEnv(env Env) (Config, error) {
	dbConfig, err := DatabaseConfigFromEnv(env)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Database: dbConfig,
	}, nil
}
