package config

type DatabaseConfig struct {
	Driver       string
	DSN          string
	MigrationDir string
}

func DatabaseConfigFromEnv(env Env) (DatabaseConfig, error) {
	return DatabaseConfig{
		Driver:       env["DATABASE_DRIVER"],
		DSN:          env["DATABASE_DSN"],
		MigrationDir: env["DATABASE_MIGRATIONS_DIR"],
	}, nil
}
