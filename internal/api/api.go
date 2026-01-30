package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/OkaniYoshiii/sqlite-go/internal/config"
	"github.com/OkaniYoshiii/sqlite-go/internal/database"
	"github.com/OkaniYoshiii/sqlite-go/internal/debug"
	"github.com/OkaniYoshiii/sqlite-go/internal/middleware"
	"github.com/OkaniYoshiii/sqlite-go/internal/repository"
	"github.com/OkaniYoshiii/sqlite-go/internal/routes"
	"github.com/go-playground/validator/v10"
)

func Run(address string, readTimeout, readHeaderTimeout, writeTimeout, idleTimeout int) error {
	env, err := config.LoadEnv()
	if err != nil {
		return err
	}

	conf, err := config.FromEnv(env)
	if err != nil {
		return err
	}

	logger, err := debug.NewLogger()

	db, err := database.Open(conf.Database.Driver, conf.Database.DSN)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	validate := validator.New()

	queries := repository.New()

	deps := routes.Dependencies{
		DB:      db,
		Queries: queries,
		Logger:  logger,
	}

	apiMiddleware := middleware.ApiMiddleware(db, queries, logger, validate)
	cspMiddleware := middleware.CSPMiddleware()
	authMiddleware := middleware.AuthMiddleware(*queries, db, &conf)
	corsMiddleware := middleware.CorsMiddleware(conf.CORS.TrustedOrigins)
	rateLimiterMiddleware := middleware.RateLimiterMiddleware(0.5)

	mux.HandleFunc("GET /api/v1/health", routes.HealthHandler(deps))
	mux.HandleFunc("GET /api/v1/posts", routes.PostsHandler(deps))
	mux.Handle("POST /api/v1/posts", middleware.Pipe(authMiddleware)(routes.PostsHandler(deps)))
	mux.HandleFunc("POST /api/v1/register", routes.RegisterHandler(deps))
	mux.HandleFunc("POST /api/v1/login", routes.LoginHandler(logger, validate, queries, db, &conf))

	server := http.Server{
		Addr:              address,
		ReadTimeout:       time.Millisecond * time.Duration(readTimeout),
		ReadHeaderTimeout: time.Millisecond * time.Duration(readHeaderTimeout),
		WriteTimeout:      time.Millisecond * time.Duration(writeTimeout),
		IdleTimeout:       time.Duration(idleTimeout),
		Handler: middleware.Pipe(
			corsMiddleware,
			rateLimiterMiddleware,
			apiMiddleware,
			cspMiddleware,
		)(mux),
		ErrorLog: logger,
	}

	defer func() {
		db.Close()
		if file, ok := logger.Writer().(*os.File); ok {
			file.Close()
		}
	}()

	fmt.Printf("Server listening on %s\n", address)
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
