package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/OkaniYoshiii/sqlite-go/internal/jwt"
)

type Config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Driver       string
	DSN          string
	MigrationDir string
}

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

func DatabaseConfigFromEnv(env Env) (DatabaseConfig, error) {
	return DatabaseConfig{
		Driver:       env.DatabaseDriver,
		DSN:          env.DatabaseDSN,
		MigrationDir: env.GooseMigrationDir,
	}, nil
}


	}


	}


}
