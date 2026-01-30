package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func Pipe(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := range len(middlewares) {
			next = middlewares[i](next)
		}

		return next
	}
}
