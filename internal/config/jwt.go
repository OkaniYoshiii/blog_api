package config

import (
	"strconv"
	"time"
)

type JWTEnv struct {
	Secret string
	TTL string
}

type JWTConfig struct {
	Secret string
	TTL time.Duration
}

func JWTConfigFromEnv(env Env) (JWTConfig, error) {
	secret := env.JWT.Secret
	ttl, err := strconv.Atoi(env.JWT.TTL)
	if err != nil {
		return JWTConfig{}, err
	}

	return JWTConfig{
		Secret: secret,
		TTL: time.Duration(ttl),
	}, nil
}
