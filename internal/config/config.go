package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type Env struct {
	DatabaseDriver string
	DatabaseDSN    string
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
		Driver: env.DatabaseDriver,
		DSN:    env.DatabaseDSN,
	}
}

func LoadEnv() (Env, error) {
	envMap, err := godotenv.Read()
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

	env.DatabaseDriver = driver
	env.DatabaseDSN = dsn

	return env, nil
}
