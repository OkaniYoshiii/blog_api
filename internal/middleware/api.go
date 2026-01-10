package middleware

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"slices"

	"github.com/OkaniYoshiii/sqlite-go/internal/repository"
	"github.com/go-playground/validator/v10"
)

func ApiMiddleware(next http.Handler, db *sql.DB, queries *repository.Queries, logger *log.Logger, validate *validator.Validate) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		apiKeys, err := queries.ListApiKeys(context.Background(), db)

		if err != nil {
			logger.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		apiKey := request.Header.Get("X-API-Key")

		if err := validate.Var(apiKey, "required"); err != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !slices.ContainsFunc(apiKeys, func(key repository.ListApiKeysRow) bool { return key.Value == apiKey }) {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
