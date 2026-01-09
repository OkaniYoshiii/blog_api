package routes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/OkaniYoshiii/sqlite-go/internal/repository"
	"github.com/go-playground/validator/v10"
)

func PostsHandler(deps Dependencies) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Accept-Post", "application/json")
		writer.Header().Set("Content-Type", "application/json")

		switch request.Method {
		case http.MethodGet:
			GetPostsHandler(deps, writer, request)
		case http.MethodPost:
			PostPostsHandler(deps, writer, request)
		}
	}
}

func GetPostsHandler(deps Dependencies, writer http.ResponseWriter, request *http.Request) {
	db := deps.DB
	queries := deps.Queries
	logger := deps.Logger

	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	posts, err := queries.ListPosts(context.Background(), db)
	if err != nil {
		logger.Println(err)

		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(posts)
	if err != nil {
		logger.Println(err)

		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Write(respBody)
}

func PostPostsHandler(deps Dependencies, writer http.ResponseWriter, request *http.Request) {
	db := deps.DB
	queries := deps.Queries
	logger := deps.Logger

	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if request.Header.Get("Content-Type") != "application/json" {
		writer.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	reqBody, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Println(err)

		writer.WriteHeader(http.StatusInternalServerError)
	}

	validate := validator.New()
	if err := validate.Var(reqBody, "json"); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	createPostParams := repository.CreatePostParams{}
	if err := json.Unmarshal(reqBody, &createPostParams); err != nil {
		if _, ok := err.(*json.SyntaxError); ok {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := validate.Var(createPostParams.Title, "required,alphanumspace"); err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := validate.Var(createPostParams.Content, "required,alphanumspace"); err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	post, err := queries.CreatePost(context.Background(), db, createPostParams)
	if err != nil {
		logger.Println(err)

		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(post)
	if err != nil {
		logger.Println(err)

		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Write(data)
}
