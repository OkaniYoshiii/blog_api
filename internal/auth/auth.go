package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/OkaniYoshiii/sqlite-go/internal/config"
	"github.com/OkaniYoshiii/sqlite-go/internal/jwt"
	"github.com/OkaniYoshiii/sqlite-go/internal/repository"
)

type authError struct {
	message string
	cause   error
}

func AuthError(message string, cause error) error {
	return &authError{message: message, cause: cause}
}

func (err *authError) Error() string {
	return fmt.Errorf("%s : %w", err.message, err.cause).Error()
}

type jwtAuth struct {
	User   repository.GetUserByIdRow
	Claims jwt.Claims
}

func UsingJWT(request *http.Request, queries repository.Queries, db repository.DBTX, conf *config.Config) (jwtAuth, error) {
	header := request.Header.Get("Authorization")
	token, ok := strings.CutPrefix(header, "Bearer ")
	if !ok {
		return jwtAuth{}, AuthError("access token not found in HTTP Authorization header", nil)
	}

	claims, err := jwt.ValidateAccessToken(token, *conf)
	if err != nil {
		return jwtAuth{}, AuthError("invalid access token", err)
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return jwtAuth{}, AuthError("access token subject claim cannot be converted to an integer", nil)
	}

	user, err := queries.GetUserById(context.Background(), db, int64(userId))
	if err != nil {
		return jwtAuth{}, AuthError("user not found", err)
	}

	return jwtAuth{User: user, Claims: claims}, nil
}
