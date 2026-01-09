package routes

import (
	"net/http"
)

func HealthHandler(deps Dependencies) http.HandlerFunc {
	db := deps.DB
	logger := deps.Logger

	return func(writer http.ResponseWriter, request *http.Request) {
		if db == nil {
			logger.Println("health check failed: database connection is nil")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := db.Ping(); err != nil {
			logger.Println("health check failed: cannot connect to database")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}
