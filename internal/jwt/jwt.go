// Package jwt creates and validate application specific tokens
//
// As a reminder, JWTs are composed of multiple claims which are :
//
// - Issuer (iss) : case sensitive string identifying the server who created it. In this case, this server created the token so this is the URL of this server.
//
// - Subject (sub) : identifier of the user, service or whatever else that will use this token. In the case of web auth, it can be an user id.
//
// - Audience (aud) : string or URL identifying the service verifying the JWT. In this case, this server is verifying the token so this is the URL of this server.
//
// - ExpiresAt (exp) : timestamp where the token should be considered invalid.
//
// - NotBefore (nbf) : timestamp indicating that the token is invalid unless current timestamp is greater than this. Can be used to create the token in advance, send it to a client but enabling it later. Must be less than "iat".
//
// - Issued At (iat) : timestamp indicating when the token has been created.
//
// - JWT ID (jti) : identifier for a particular token. Can be use to revoke the token if it has been compromised. Can also be used for refresh tokens.
//
// Issuer and Audience can be different for example when using Google / Facebook login.
// You will authenticate with a Facebook/Google server which will send a JWT back.
// This JWT is then used on your API to access ressources.
// In this case, the issuer is Google/Facebook and the audience is your API.
//
// Some types of token will use some claims and not other. This package aims
// to generate and validate specific tokens like access tokens, refresh tokens etc ...
// in the context of this application.
package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/OkaniYoshiii/sqlite-go/internal/config"
	"github.com/OkaniYoshiii/sqlite-go/internal/repository"
	"github.com/golang-jwt/jwt/v4"
)

// Minimum number of bits the JWT secret must have
// See : https://cheatsheetseries.owasp.org/cheatsheets/JSON_Web_Token_for_Java_Cheat_Sheet.html#weak-token-secret
const SecretMinStrength = 256
const BitsInByte = 8

var SigningMethod = jwt.SigningMethodHS256

type Claims struct {
	jwt.RegisteredClaims
}

func ValidateSecret(secret []byte) error {
	if len(secret)*BitsInByte <= SecretMinStrength {
		return fmt.Errorf("secret must be %d bits long", SecretMinStrength)
	}

	return nil
}

// Creates a new signed JWT Token using an HMAC based signing method and a secret
// Must be used to authenticate a user
// Can be validated by using [ValidateAccessToken]
func NewAccessToken(conf config.Config, user repository.User) (string, error) {
	ttl := conf.JWT.TTL

	issuer := accessTokenIssuer(conf)
	audience := accessTokenAudience(conf)
	subject := strconv.Itoa(int(user.ID))
	expiresAt := time.Now().Add(ttl)
	notBefore := time.Now()
	issuedAt := time.Now()
	id := ""

	secret := []byte(conf.JWT.Secret)

	if err := ValidateSecret(secret); err != nil {
		return "", err
	}

	claims := Claims{
		jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   subject,
			Audience:  jwt.ClaimStrings(audience),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(notBefore),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ID:        id,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signed, nil
}

// Validate an access token generated with [NewAccessToken]
// Since [jwt.RegisteredClaims.Valid] does not require that "exp", "iat" and "nbf" claims are present
// manual validation is made for this type of token.
func ValidateAccessToken(tokenStr string, conf config.Config) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		return conf.JWT.Secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return Claims{}, err
	}

	claims := token.Claims.(*Claims)

	if !claims.VerifyExpiresAt(time.Now(), true) {
		return Claims{}, jwt.ErrTokenExpired
	}

	if !claims.VerifyIssuedAt(time.Now(), true) {
		return Claims{}, jwt.ErrTokenUsedBeforeIssued
	}

	if !claims.VerifyNotBefore(time.Now(), true) {
		return Claims{}, jwt.ErrTokenNotValidYet
	}

	if !claims.VerifyIssuer(accessTokenIssuer(conf), true) {
		return Claims{}, jwt.ErrTokenInvalidIssuer
	}

	for _, aud := range accessTokenAudience(conf) {
		if !claims.VerifyAudience(aud, true) {
			return Claims{}, jwt.ErrTokenInvalidAudience
		}
	}

	return claims, nil
}

func accessTokenIssuer(conf config.Config) string {
	return conf.Server.Host
}

func accessTokenAudience(conf config.Config) []string {
	return []string{conf.Server.Host}
}
