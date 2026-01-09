package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/OkaniYoshiii/sqlite-go/internal/repository"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(deps Dependencies) http.HandlerFunc {
	logger := deps.Logger
	queries := deps.Queries
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if request.Header.Get("Content-Type") != "application/json" {
			writer.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		createUserParams := repository.CreateUserParams{}

		decoder := json.NewDecoder(request.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&createUserParams); err != nil {
			writer.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		email := createUserParams.Email
		password := createUserParams.Password

		validate := validator.New()

		if err := validate.Var(email, "required,email"); err != nil {
			writer.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		// bcrypt.GenerateFromPassword does not accept password longer than 72 bytes
		if err := validate.Var(password, "required,min=8,max=72"); err != nil {
			writer.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		createUserParams.Password = string(hash)

		user, err := queries.CreateUser(context.Background(), createUserParams)
		if err != nil {
			logger.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(user)
		if err != nil {
			logger.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusCreated)
		writer.Write(body)
	}
}
