package config

type DatabaseEnv struct {
	Driver    string
	DSN       string
	MigrationsDir string
}

type DatabaseConfig struct {
	Driver       string
	DSN          string
	MigrationDir string
}

func DatabaseConfigFromEnv(env Env) (DatabaseConfig, error) {
	return DatabaseConfig{
		Driver:       env.Database.Driver,
		DSN:          env.Database.DSN,
		MigrationDir: env.Database.MigrationsDir,
	}, nil
}
