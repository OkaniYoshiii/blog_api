package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Driver       string
	DSN          string
	MigrationDir string
}

type Env struct {
	DatabaseDriver    string
	DatabaseDSN       string
	GooseMigrationDir string
}

func Load() (Config, error) {
	env, err := LoadEnv()
	if err != nil {
		return Config{}, err
	}

	dbConfig := LoadDatabaseConfig(env)

	return Config{
		Database: dbConfig,
	}, nil
}

func LoadDatabaseConfig(env Env) DatabaseConfig {
	return DatabaseConfig{
		Driver:       env.DatabaseDriver,
		DSN:          env.DatabaseDSN,
		MigrationDir: env.GooseMigrationDir,
	}
}

func LoadEnv(filenames ...string) (Env, error) {
	envMap, err := godotenv.Read(filenames...)
	if err != nil {
		return Env{}, err
	}

	env := Env{}
	errFmt := "Environment variable %s must be set and must not be empty."

	driver, ok := envMap["DATABASE_DRIVER"]
	if !ok || len(driver) == 0 {
		return Env{}, fmt.Errorf(errFmt, "DATABASE_DRIVER")
	}

	dsn, ok := envMap["DATABASE_DSN"]
	if !ok || len(dsn) == 0 {
		return Env{}, fmt.Errorf(errFmt, "DATABASE_DSN")
	}

	gooseMigrationDir, ok := envMap["GOOSE_MIGRATION_DIR"]
	if !ok || len(gooseMigrationDir) == 0 {
		return Env{}, fmt.Errorf(errFmt, "GOOSE_MIGRATION_DIR")
	}

	env.DatabaseDriver = driver
	env.DatabaseDSN = dsn
	env.GooseMigrationDir = gooseMigrationDir

	return env, nil
}
