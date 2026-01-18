package config

const (
	DatabaseDSN Key = "DATABASE_DSN"
	DatabaseDriver Key = "DATABASE_DRIVER"
	DatabaseMigrationsDir Key = "DATABASE_MIGRATIONS_DIR"
)

type DatabaseConfig struct {
	Driver       string
	DSN          string
	MigrationsDir string
}

type DatabaseEnv struct {
	Driver       string
	DSN          string
	MigrationsDir string
}

func DatabaseConfigFromEnv(env Env) (DatabaseConfig, error) {
	if err := ValidateDefined(DatabaseDSN, env.Database.DSN); err != nil {
		return DatabaseConfig{}, err
	}

	if err := ValidateDefined(DatabaseDriver, env.Database.DSN); err != nil {
		return DatabaseConfig{}, err
	}

	if err := ValidateDefined(DatabaseMigrationsDir, env.Database.MigrationsDir); err != nil {
		return DatabaseConfig{}, err
	}

	return DatabaseConfig{
		Driver:       env.Database.Driver,
		DSN:          env.Database.DSN,
		MigrationsDir: env.Database.MigrationsDir,
	}, nil
}
