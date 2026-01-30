package middleware

import (
	"net/http"

	"github.com/OkaniYoshiii/sqlite-go/pkg/csp"
)

func CSPMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			directives := csp.Strict()

			writer.Header().Set("Content-Security-Policy", directives.String())

			next.ServeHTTP(writer, request)
		})
	}
}
