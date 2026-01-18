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

func New(issuer string, subject string, audience jwt.ClaimStrings, expiresAt *jwt.NumericDate, notBefore *jwt.NumericDate, id string, secret []byte) (string, error) {
	if err := ValidateSecret(secret); err != nil {
		return "", err
	}

	claims := Claims{
		jwt.RegisteredClaims{
			Issuer: issuer,
			Subject: subject,
			Audience: audience,
			ExpiresAt: expiresAt,
			NotBefore: notBefore,
			IssuedAt: jwt.NewNumericDate(time.Now()),
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
	if len(secret) * BitsInByte != SecretMinStrength {
		return fmt.Errorf("secret must be %d bits long", SecretMinStrength)
	}

	return nil
}
