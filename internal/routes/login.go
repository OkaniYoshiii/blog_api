package routes

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"

	"github.com/OkaniYoshiii/sqlite-go/internal/config"
	"github.com/OkaniYoshiii/sqlite-go/internal/jwt"
	"github.com/OkaniYoshiii/sqlite-go/internal/repository"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(
	logger *log.Logger,
	validate *validator.Validate,
	queries *repository.Queries,
	db repository.DBTX,
	conf *config.Config,
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			PostLoginHandler(writer, request, logger, validate, queries, db, conf)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func PostLoginHandler(
	writer http.ResponseWriter,
	request *http.Request,
	logger *log.Logger,
	validate *validator.Validate,
	queries *repository.Queries,
	db repository.DBTX,
	conf *config.Config,
) {
	writer.Header().Set("Content-Type", "application/json")

	if !slices.Contains(request.Header["Accept"], "application/json") {
		writer.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if !slices.Contains(request.Header["Content-Type"], "application/json") {
		writer.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	credentials := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	decoder := json.NewDecoder(request.Body)

	if err := decoder.Decode(&credentials); err != nil {
		switch err := err.(type) {
		case *json.SyntaxError, *json.UnsupportedValueError:
			writer.WriteHeader(http.StatusBadRequest)
			return
		default:
			logger.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if err := validate.Var(credentials.Email, "required,email"); err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := validate.Var(credentials.Password, "required,min=8"); err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	user, err := queries.GetUserByEmail(context.Background(), db, credentials.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		logger.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		logger.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := jwt.NewAccessToken(*conf, user)
	if err != nil {
		logger.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody := struct {
		AccessToken string `json:"access_token"`
	}{token}

	encoder := json.NewEncoder(writer)
	if err := encoder.Encode(responseBody); err != nil {
		logger.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
