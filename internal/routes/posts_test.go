package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OkaniYoshiii/sqlite-go/internal/config"
	"github.com/OkaniYoshiii/sqlite-go/internal/database"
	"github.com/OkaniYoshiii/sqlite-go/internal/repository"
	"github.com/pressly/goose/v3"
)

type RouteTest struct {
	name     string
	data     TestRequest
	expected TestResponse
}

type TestResponse struct {
	body       string
	statusCode int
}

type TestRequest struct {
	body   string
	method string
	url    string
	header http.Header
}

func (resp *TestResponse) AssertEqual(t *testing.T, other TestResponse) {
	if resp.statusCode != other.statusCode {
		t.Errorf("mismatched status code: expected %d, got %d", resp.statusCode, other.statusCode)
	}

	if resp.body != other.body {
		t.Errorf("mismatched content: expected %s, got: %s", resp.body, other.body)
	}
}

var (
	env config.Env
)

func TestMain(m *testing.M) {
	var err error
	env, err = config.LoadEnv("../../.env.test")
	if err != nil {
		log.Fatal(err)
	}

	// Empêche les logs de Goose
	buffer := strings.Builder{}
	log.SetOutput(&buffer)

	os.Exit(m.Run())
}

func Setup(t *testing.T) Dependencies {
	db, err := database.Open(env["DATABASE_DRIVER"], env["DATABASE_DSN"])
	if err != nil {
		log.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	logger := log.New(&strings.Builder{}, "", 0)
	queries := repository.New()
	deps := Dependencies{
		DB:      db,
		Queries: queries,
		Logger:  logger,
	}

	migrationDir := filepath.Join("../../", env["DATABASE_MIGRATIONS_DIR"])

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db, migrationDir); err != nil {
		log.Fatal(err)
	}

	return deps
}

func MustMarshal(t *testing.T, data any) []byte {
	result, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	return result
}

func NewCreatePostParams() repository.CreatePostParams {
	return repository.CreatePostParams{
		Title:   "Test",
		Content: "Test content",
	}
}

func PostFromCreateParams(id int64, params repository.CreatePostParams) repository.Post {
	return repository.Post{
		ID:      id,
		Title:   params.Title,
		Content: params.Content,
	}
}

func TestPostsRead(t *testing.T) {
	tests := [1]RouteTest{}

	createPostParams := NewCreatePostParams()

	posts := []repository.Post{
		PostFromCreateParams(1, createPostParams),
	}

	tests[0].name = "Successfull read"
	tests[0].data.body = ""
	tests[0].data.method = "GET"
	tests[0].data.url = "/api/posts"
	tests[0].expected.body = string(MustMarshal(t, posts))
	tests[0].expected.statusCode = http.StatusOK

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			deps := Setup(t)
			queries := deps.Queries
			db := deps.DB

			_, err := queries.CreatePost(context.Background(), db, createPostParams)
			if err != nil {
				t.Fatal(err)
			}

			response := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", "/api/posts", nil)

			handler := PostsHandler(deps)
			handler(response, request)

			got := TestResponse{
				statusCode: response.Result().StatusCode,
				body:       response.Body.String(),
			}

			test.expected.AssertEqual(t, got)
		})
	}
}

func TestPostsCreate(t *testing.T) {
	deps := Setup(t)

	createPostParams := NewCreatePostParams()
	post := PostFromCreateParams(1, createPostParams)

	tests := [6]RouteTest{}

	tests[0].name = "Successfull create"
	tests[0].data.body = string(MustMarshal(t, createPostParams))
	tests[0].data.method = "POST"
	tests[0].data.url = "/api/posts"
	tests[0].data.header = http.Header{"Content-Type": []string{"application/json"}}
	tests[0].expected.body = string(MustMarshal(t, post))
	tests[0].expected.statusCode = http.StatusCreated

	tests[1].name = "Invalid body format"
	tests[1].data.body = "azeadavvdaad"
	tests[1].data.method = "POST"
	tests[1].data.header = http.Header{"Content-Type": []string{"application/json"}}
	tests[1].data.url = "/api/posts"
	tests[1].expected.body = ""
	tests[1].expected.statusCode = http.StatusBadRequest

	tests[2].name = "Invalid JSON Post"
	tests[2].data.body = `{"title": "Correct title", "undefined_key": "Not valid"}`
	tests[2].data.method = "POST"
	tests[2].data.header = http.Header{"Content-Type": []string{"application/json"}}
	tests[2].data.url = "/api/posts"
	tests[2].expected.body = ""
	tests[2].expected.statusCode = http.StatusUnprocessableEntity

	tests[3].name = "Empty title"
	tests[3].data.body = `{"title": "", "content": "Post content"}`
	tests[3].data.method = "POST"
	tests[3].data.header = http.Header{"Content-Type": []string{"application/json"}}
	tests[3].data.url = "/api/posts"
	tests[3].expected.body = ""
	tests[3].expected.statusCode = http.StatusUnprocessableEntity

	tests[4].name = "Empty content"
	tests[4].data.body = `{"title": "Title", "content": ""}`
	tests[4].data.method = "POST"
	tests[4].data.header = http.Header{"Content-Type": []string{"application/json"}}
	tests[4].data.url = "/api/posts"
	tests[4].expected.body = ""
	tests[4].expected.statusCode = http.StatusUnprocessableEntity

	tests[5].name = "Invalid Content-Type"
	tests[5].data.body = string(MustMarshal(t, createPostParams))
	tests[5].data.method = "POST"
	tests[5].data.header = http.Header{"Content-Type": []string{"text/plain"}}
	tests[5].data.url = "/api/posts"
	tests[5].expected.body = ""
	tests[5].expected.statusCode = http.StatusUnsupportedMediaType

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodPost, "/api/posts", strings.NewReader(test.data.body))
			for key, values := range test.data.header {
				for _, value := range values {
					request.Header.Set(key, value)
				}
			}

			handler := PostsHandler(deps)
			handler(response, request)

			got := TestResponse{
				body:       response.Body.String(),
				statusCode: response.Result().StatusCode,
			}

			test.expected.AssertEqual(t, got)
		})
	}
}
