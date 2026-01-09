package routes

import (
	"database/sql"
	"log"

	"github.com/OkaniYoshiii/sqlite-go/internal/repository"
)

type Dependencies struct {
	DB      *sql.DB
	Queries *repository.Queries
	Logger  *log.Logger
}
