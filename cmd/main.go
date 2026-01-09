package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/OkaniYoshiii/sqlite-go/internal/config"
	"github.com/OkaniYoshiii/sqlite-go/internal/database"
	"github.com/OkaniYoshiii/sqlite-go/internal/debug"
	"github.com/OkaniYoshiii/sqlite-go/internal/repository"
	"github.com/OkaniYoshiii/sqlite-go/internal/routes"
)

var address = flag.String("address", "127.0.0.1:8000", "Specifies the TCP address for the server to listen on, in the form “host:port”. ")
var readTimeout = flag.Int("readtimeout", 10000, "The maximum duration in milliseconds for reading the entire request, including the body.")
var readHeaderTimeout = flag.Int("readheadertimeout", 2000, "The maximum duration in milliseconds for reading the headers.")
var writeTimeout = flag.Int("writetimeout", 3000, "The maximum duration in milliseconds before timing out writes of the response.")
var idleTimeout = time.Millisecond * 100

func main() {
	flag.Parse()

	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := debug.NewLogger()

	db, err := database.Open(config.Database.Driver, config.Database.DSN)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	queries := repository.New()

	deps := routes.Dependencies{
		DB:      db,
		Queries: queries,
		Logger:  logger,
	}

	mux.HandleFunc("GET /api/v1/health", routes.HealthHandler(deps))
	mux.HandleFunc("GET /api/v1/posts", routes.PostsHandler(deps))
	mux.HandleFunc("POST /api/v1/posts", routes.PostsHandler(deps))
	mux.HandleFunc("POST /api/v1/register", routes.RegisterHandler(deps))

	server := http.Server{
		Addr:              *address,
		ReadTimeout:       time.Millisecond * time.Duration(*readTimeout),
		ReadHeaderTimeout: time.Millisecond * time.Duration(*readHeaderTimeout),
		WriteTimeout:      time.Millisecond * time.Duration(*writeTimeout),
		IdleTimeout:       idleTimeout,
		Handler:           mux,
		ErrorLog:          logger,
	}

	defer func() {
		db.Close()
		if file, ok := logger.Writer().(*os.File); ok {
			file.Close()
		}
	}()

	fmt.Printf("Server listening on %s\n", *address)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
