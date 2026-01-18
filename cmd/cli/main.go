package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/OkaniYoshiii/sqlite-go/internal/config"
	"github.com/OkaniYoshiii/sqlite-go/internal/database"
	"github.com/OkaniYoshiii/sqlite-go/internal/jwt"
	"github.com/OkaniYoshiii/sqlite-go/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

// Outil en ligne de commande pour gérer les clés API
func main() {
	args := os.Args

	if err := validate.VarWithKey("command_name", args, "required,min=2"); err != nil {
		log.Fatal(err)
	}

	cmd := args[1]

	var err error
	switch cmd {
	case "apikey:generate":
		err = ApiCommand(args)
	case "jwt:generate-secret":
	  err = JWTGenerateSecret()
	default:
		err = fmt.Errorf("%s: command not found", cmd)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func ApiCommand(args []string) error {
	applications := []string{"web_backend", "web_frontend"}
	if err := validate.Var(args, "len=3"); err != nil {
		return fmt.Errorf("missing argument \"application\". Possible values are \"%s\"", applications)
	}

	application := args[2]

	if !slices.Contains(applications, application) {
		return fmt.Errorf("%s is not a valid application. Possible values are %v", application, applications)
	}

	env, err := config.LoadEnv()
	if err != nil {
		return err
	}

	cfg, err := config.FromEnv(env)
	if err != nil {
		return err
	}

	db, err := database.Open(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		return err
	}

	queries := repository.New()
	apiKey := uuid.NewString()

	createApiKeyParams := repository.CreateApiKeyParams{
		Value:       apiKey,
		Application: application,
	}

	if err := queries.CreateApiKey(context.Background(), db, createApiKeyParams); err != nil {
		return err
	}

	fmt.Printf("New API Key created : %s\n", apiKey)

	return nil
}

func JWTGenerateSecret() error {
	b := [jwt.SecretMinStrength / 8]byte{}
	_, err := rand.Read(b[:])
	if err != nil {
		return err
	}

	secret := base64.StdEncoding.EncodeToString(b[:])

	fmt.Printf("Generated %d-bit key: %s\n", jwt.SecretMinStrength, secret)

	return nil
}
