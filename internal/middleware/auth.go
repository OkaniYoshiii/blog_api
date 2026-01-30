package middleware

import (
	"context"
	"net/http"

	"github.com/OkaniYoshiii/sqlite-go/internal/auth"
	"github.com/OkaniYoshiii/sqlite-go/internal/config"
	"github.com/OkaniYoshiii/sqlite-go/internal/repository"
)

func AuthMiddleware(queries repository.Queries, db repository.DBTX, conf *config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			authentication, err := auth.UsingJWT(request, queries, db, conf)
			if err != nil {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(request.Context(), "auth", authentication)

			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}
