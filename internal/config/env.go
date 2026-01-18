package config

import "github.com/joho/godotenv"

type Env struct {
	DatabaseDriver    string
	DatabaseDSN       string
	GooseMigrationDir string
}

func LoadEnv(filenames ...string) (Env, error) {
	envMap, err := godotenv.Read(filenames...)
	if err != nil {
		return Env{}, err
	}

	env := Env{}
	env.DatabaseDriver = envMap["DATABASE_DRIVER"]
	env.DatabaseDSN = envMap["DATABASE_DSN"]
	env.GooseMigrationDir = envMap["GOOSE_MIGRATIONS_DIR"]

	return env, nil
}
