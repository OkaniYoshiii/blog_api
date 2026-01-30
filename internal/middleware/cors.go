package middleware

import "net/http"

func CorsMiddleware(origins []string) Middleware {
	cors := http.NewCrossOriginProtection()

	for _, origin := range origins {
		cors.AddTrustedOrigin(origin)
	}

	return cors.Handler
}
