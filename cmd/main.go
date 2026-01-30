package main

import (
	"log"

	"github.com/OkaniYoshiii/sqlite-go/internal/api"
)

func main() {
	if err := api.Run(); err != nil {
		log.Fatal(err)
	}
}
