package config

import "github.com/joho/godotenv"

type Env = map[string]string

func LoadEnv(filenames ...string) (map[string]string, error) {
	env, err := godotenv.Read(filenames...)
	if err != nil {
		return map[string]string{}, err
	}

	return env, nil
}
