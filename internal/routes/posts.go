package routes

import (
	"context"
	"encoding/json"
	"net/http"
)

func PostsHandler(deps Dependencies) http.HandlerFunc {
	queries := deps.Queries
	logger := deps.Logger

	return func(writer http.ResponseWriter, request *http.Request) {
		posts, err := queries.ListPosts(context.Background())
		writer.Header().Set("Content-Type", "application/json")
		if err != nil {
			logger.Println(err)

			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte{})
			return
		}

		body, err := json.Marshal(posts)
		if err != nil {
			logger.Println(err)

			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte{})
			return
		}

		writer.Write(body)
	}
}
