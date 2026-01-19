package jwt

import (
	"fmt"
	"time"

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

// Creates a new signed JWT Token using an HMAC based signing method and a secret
//
// JWTs are composed of multiple claims which are :
// - Issuer (iss) : case sensitive string identifying the server who created it. In this case, this server created the token so this is the URL of this server.
// - Subject (sub) : identifier of the user, service or whatever else that will use this token. In the case of web auth, it can be an user id.
// - Audience (aud) : string or URL identifying the service verifying the JWT. In this case, this server is verifying the token so this is the URL of this server.
// - ExpiresAt (exp) : timestamp where the token should be considered invalid.
// - NotBefore (nbf) : timestamp indicating that the token is invalid unless current timestamp is greater than this. Can be used to create the token in advance, send it to a client but enabling it later. Must be less than "iat".
// - Issued At (iat) : timestamp indicating when the token has been created.
// - JWT ID (jti) : identifier for a particular token. Can be use to revoke the token if it has been compromised. Can also be used for refresh tokens.
//
// Issued at is automatically assigned with the current timestamp.
// Other fields are optionnal and will not be included in the JWT if empty.
//
// Issuer and Audience can be different for example when using Google / Facebook login.
// You will authenticate with a Facebook/Google server which will send a JWT back.
// This JWT is then used on your API to access ressources.
// In this case, the issuer is Google/Facebook and the audience is your API.
func New(issuer string, subject string, audience []string, expiresAt, notBefore, issuedAt time.Time, id string, secret []byte) (string, error) {
	if err := ValidateSecret(secret); err != nil {
		return "", err
	}

	claims := Claims{
		jwt.RegisteredClaims{
			Issuer: issuer,
			Subject: subject,
			Audience: jwt.ClaimStrings(audience),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(notBefore),
			IssuedAt: jwt.NewNumericDate(issuedAt),
			ID: id,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ValidateSecret(secret []byte) error {
	if len(secret) * BitsInByte <= SecretMinStrength {
		return fmt.Errorf("secret must be %d bits long", SecretMinStrength)
	}

	return nil
}
