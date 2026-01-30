package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/OkaniYoshiii/sqlite-go/internal/config"
	"github.com/OkaniYoshiii/sqlite-go/internal/jwt"
	"github.com/OkaniYoshiii/sqlite-go/internal/repository"
)

type Service struct {
	Conf config.Config
}

func (service *Service) FromRequest(request *http.Request) (repository.User, error) {
	header := request.Header.Get("Authorization")
	token, ok := strings.CutPrefix(header, "Bearer ")
	if !ok {
		return repository.User{}, fmt.Errorf("")
	}

	claims, err := jwt.ValidateAccessToken(token, service.Conf)
	if err != nil {
		return repository.User{}, err
	}

	claims
}
