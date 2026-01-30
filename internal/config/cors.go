package config

import "strings"

type CORSConfig struct {
	TrustedOrigins []string
}

func CORSConfigFromEnv(env Env) (CORSConfig, error) {
	trustedOrigins := strings.Split(env["CORS_TRUSTED_ORIGINS"], ",")

	return CORSConfig{
		TrustedOrigins: trustedOrigins,
	}, nil
}
