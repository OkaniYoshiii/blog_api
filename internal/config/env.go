package config

import (
	"github.com/joho/godotenv"
)

type Key = string

type Env struct {
	Database DatabaseEnv
	JWT JWTEnv
}

func LoadEnv(filenames ...string) (Env, error) {
	envMap, err := godotenv.Read(filenames...);
	if err != nil {
		return Env{}, nil
	}

	env := Env{}

	env.JWT.Secret = envMap[JWTSecret]
	env.JWT.TTL = envMap[JWTTTL]

	env.Database.DSN = envMap[DatabaseDSN]
	env.Database.Driver = envMap[DatabaseDriver]
	env.Database.MigrationsDir = envMap[DatabaseMigrationsDir]

	return env, nil
}
