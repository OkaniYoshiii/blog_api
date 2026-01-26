package routes

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/OkaniYoshiii/sqlite-go/internal/config"
	"github.com/OkaniYoshiii/sqlite-go/internal/database"
	"github.com/OkaniYoshiii/sqlite-go/internal/repository"
	"github.com/go-playground/validator/v10"
)

func TestLogin(t *testing.T) {
	tests := [1]struct {
		Name    string
		Request struct {
			Method string
			Body   io.Reader
			Header http.Header
		}
		Expected struct {
			StatusCode int
			Header     http.Header
		}
	}{}

	env, err := config.LoadEnv("../../.env.test")
	if err != nil {
		t.Fatal(err)
	}

	conf, err := config.FromEnv(env)
	if err != nil {
		t.Fatal(err)
	}

	logger := log.New(&strings.Builder{}, log.Default().Prefix(), log.Default().Flags())
	queries := repository.New()
	validate := validator.New(validator.WithRequiredStructEnabled())
	db, err := database.Open(conf.Database.Driver, conf.Database.DSN)
	if err != nil {
		t.Fatal(err)
	}

	createUserParams := repository.CreateUserParams{
		Email:    "user@mail.com",
		Password: "password_longer_than_8_caracters",
	}

	queries.CreateUser(context.Background(), db, createUserParams)
	tests[0].Name = "Successfull test"
	tests[0].Request.Method = "POST"
	tests[0].Request.Body = strings.NewReader(`{"email": "", "password": ""}`)
	tests[0].Request.Header = http.Header{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"application/json"},
	}
	tests[0].Expected.StatusCode = http.StatusOK
	tests[0].Expected.Header = http.Header{
		"Content-Type": []string{"application/json"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest(test.Request.Method, "/api/v1/login", test.Request.Body)

			handler := LoginHandler(logger, validate, queries, db, &conf)
			handler(recorder, request)

			response := recorder.Result()

			if response.StatusCode != test.Expected.StatusCode {
				t.Errorf("unexpected response status code : expected %d, got %d", test.Expected.StatusCode, response.StatusCode)
			}

			for header := range test.Expected.Header {
				expected := test.Expected.Header.Get(header)
				got := response.Header.Get(header)

				if expected != got {
					t.Errorf("unexpected response header %q : expected %#v, got %#v", header, expected, got)
				}
			}
		})
	}
}
