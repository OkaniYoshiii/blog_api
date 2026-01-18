package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func LoginHandler(deps Dependencies, validate *validator.Validate) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
			case http.MethodOptions:
				OptionsLoginHandler(writer)
			case http.MethodPost:
				PostLoginHandler(deps, validate, writer, request)
			default:
				writer.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func OptionsLoginHandler(writer http.ResponseWriter) {
	allowedMethods := [...]string{http.MethodOptions, http.MethodPost}
	writer.Header().Set("Accept", strings.Join(allowedMethods[:], ", "))
}

func PostLoginHandler(
	deps Dependencies,
	validate *validator.Validate,
	writer http.ResponseWriter,
	request *http.Request,
) {
	logger := deps.Logger
	db := deps.DB
	queries := deps.Queries

	if request.Header.Get("Content-Type") != "application/json" {
		writer.Header().Set("Accept-Post", "application/json")
		writer.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	if request.Header.Get("Accept") != "application/json" {
		writer.WriteHeader(http.StatusNotAcceptable)
		return
	}

	credentials := struct{
		Email string
		Password string
	}{}

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&credentials); err != nil {
		switch err := err.(type) {
			case *json.SyntaxError:
				writer.WriteHeader(http.StatusBadRequest)
			case *json.UnmarshalTypeError:
				writer.WriteHeader(http.StatusUnprocessableEntity)
			default:
				logger.Println(err)
				writer.WriteHeader(http.StatusInternalServerError)
		}
	}

	email := credentials.Email
	password := credentials.Password

	if err := validate.Var(email, "required,email"); err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := validate.Var(password, "required,min=8"); err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	_, err := queries.GetUserByEmail(context.Background(), db, email)
	if err != nil {
		logger.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
