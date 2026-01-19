package config

import "github.com/joho/godotenv"

type Env struct {
	Database DatabaseEnv
}

func LoadEnv(filenames ...string) (Env, error) {
	envMap, err := godotenv.Read(filenames...)
	if err != nil {
		return Env{}, err
	}

	env := Env{}

	env.Database.Driver = envMap["DATABASE_DRIVER"]
	env.Database.DSN = envMap["DATABASE_DSN"]
	env.Database.MigrationsDir = envMap["GOOSE_MIGRATIONS_DIR"]

	return env, nil
}
