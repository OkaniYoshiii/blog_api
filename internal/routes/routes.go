package routes

import (
	"log"

	"github.com/OkaniYoshiii/sqlite-go/internal/repository"
)

type Dependencies struct {
	Queries *repository.Queries
	Logger  *log.Logger
}
