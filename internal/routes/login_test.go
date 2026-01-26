package routes

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/OkaniYoshiii/sqlite-go/internal/config"
	"github.com/OkaniYoshiii/sqlite-go/internal/database"
	"github.com/OkaniYoshiii/sqlite-go/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	env, err := config.LoadEnv("../../.env.test")
	if err != nil {
		t.Fatal(err)
	}

	conf := config.Config{
		Database: config.DatabaseConfig{
			Driver: env["DATABASE_DRIVER"],
			DSN:    env["DATABASE_DSN"],
		},
		Server: config.ServerConfig{
			Host: "example.com",
			Port: 80,
		},
		JWT: config.JWTConfig{
			Secret: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			TTL:    time.Duration(100) * time.Second,
		},
	}

	logsBuffer := strings.Builder{}
	logger := log.New(&logsBuffer, log.Default().Prefix(), log.Default().Flags())
	queries := repository.New()
	validate := validator.New(validator.WithRequiredStructEnabled())
	db, err := database.Open(conf.Database.Driver, conf.Database.DSN)
	if err != nil {
		t.Fatal(err)
	}

	migrationDir := filepath.Join("../../", env["DATABASE_MIGRATIONS_DIR"])

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db, migrationDir); err != nil {
		log.Fatal(err)
	}

	email := "user@mail.com"
	password := "password_longer_than_8_caracters"

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := queries.CreateUser(context.Background(), db, repository.CreateUserParams{Email: email, Password: string(hashedPwd)}); err != nil {
		t.Fatal(err)
	}

	tests := [4]struct {
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

	tests[0].Name = "Successfull test"
	tests[0].Request.Method = "POST"
	tests[0].Request.Body = strings.NewReader(fmt.Sprintf(`{"email": %q, "password": %q}`, email, password))
	tests[0].Request.Header = http.Header{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"application/json"},
	}
	tests[0].Expected.StatusCode = http.StatusOK
	tests[0].Expected.Header = http.Header{
		"Content-Type": []string{"application/json"},
	}

	tests[1].Name = "Incorrect HTTP method"
	tests[1].Request.Method = "GET"
	tests[1].Request.Body = strings.NewReader(fmt.Sprintf(`{"email": %q, "password": %q}`, email, password))
	tests[1].Request.Header = http.Header{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"application/json"},
	}
	tests[1].Expected.StatusCode = http.StatusMethodNotAllowed
	tests[1].Expected.Header = http.Header{}

	tests[2].Name = "Incorrect email"
	tests[2].Request.Method = http.MethodPost
	tests[2].Request.Body = strings.NewReader(fmt.Sprintf(`{"email": %q, "password": %q}`, "unknown@mail.com", password))
	tests[2].Request.Header = http.Header{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"application/json"},
	}
	tests[2].Expected.StatusCode = http.StatusUnauthorized
	tests[2].Expected.Header = http.Header{
		"Content-Type": []string{"application/json"},
	}

	tests[3].Name = "Incorrect password"
	tests[3].Request.Method = http.MethodPost
	tests[3].Request.Body = strings.NewReader(fmt.Sprintf(`{"email": %q, "password": %q}`, email, "wrong_password"))
	tests[3].Request.Header = http.Header{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"application/json"},
	}
	tests[3].Expected.StatusCode = http.StatusUnauthorized
	tests[3].Expected.Header = http.Header{
		"Content-Type": []string{"application/json"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest(test.Request.Method, "/api/v1/login", test.Request.Body)
			request.Header = test.Request.Header

			handler := LoginHandler(logger, validate, queries, db, &conf)
			handler(recorder, request)

			response := recorder.Result()

			if response.StatusCode != test.Expected.StatusCode {
				t.Log(logsBuffer.String())
				t.Errorf("unexpected response status code : expected %d, got %d", test.Expected.StatusCode, response.StatusCode)
			}

			for header := range test.Expected.Header {
				expected := test.Expected.Header.Get(header)
				got := response.Header.Get(header)

				if expected != got {
					t.Log(logsBuffer.String())
					t.Errorf("unexpected response header %q : expected %#v, got %#v", header, expected, got)
				}
			}
		})
	}
}
