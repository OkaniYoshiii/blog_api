package config

import (
	"fmt"
	"time"

	"github.com/OkaniYoshiii/sqlite-go/internal/jwt"
)

const (
	JWTSecret Key = "JWT_SECRET"
	JWTTTL Key = "JWT_TTL"
)

type JWTConfig struct {
	Secret string
	TTL time.Duration
}

type JWTEnv struct {
	Secret string
	TTL string
}

func JWTConfigFromEnv(env Env) (JWTConfig, error) {
	if err := ValidateDefined(JWTSecret, env.JWT.Secret); err != nil {
		return JWTConfig{}, err
	}

	if err := ValidateDefined(JWTTTL, env.JWT.TTL); err != nil {
		return JWTConfig{}, err
	}

	secret := env.JWT.Secret
	if err := jwt.ValidateSecret([]byte(secret)); err != nil {
		return JWTConfig{}, &InvalidEnvError{"JWT_SECRET", err}
	}

	ttl, err := ParseInt(JWTTTL, env.JWT.TTL)
	if err != nil {
		return JWTConfig{}, err
	}

	if ttl <= 0 {
		err := fmt.Errorf("%q is not a positive integer value", env.JWT.TTL)
		return JWTConfig{}, &InvalidEnvError{"JWT_TTL", err}
	}

	return JWTConfig{
		Secret: secret,
		TTL: time.Duration(int64(ttl)) * time.Millisecond,
	}, nil
}
