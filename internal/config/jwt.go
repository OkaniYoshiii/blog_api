package config

import (
	"strconv"
	"time"
)

type JWTConfig struct {
	Secret string
	TTL time.Duration
}

func JWTConfigFromEnv(env Env) (JWTConfig, error) {
	secret := env["JWT_SECRET"]
	ttl, err := strconv.Atoi(env["JWT_TTL"])
	if err != nil {
		return JWTConfig{}, err
	}

	return JWTConfig{
		Secret: secret,
		TTL: time.Duration(ttl),
	}, nil
}
