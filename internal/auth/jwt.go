package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()
	clms := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userId.String(),
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, clms)
	return tkn.SignedString([]byte(tokenSecret))

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	regClaims := &jwt.RegisteredClaims{}
	parsedToken, errParse := jwt.ParseWithClaims(
		tokenString,
		regClaims,
		func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("signing method is not HS256")
			}
			return []byte(tokenSecret), nil
		})
	if errParse != nil {
		return uuid.Nil, errParse
	}
	if !parsedToken.Valid {
		return uuid.Nil, errors.New("token is invalid")
	}
	claims, claimsOk := parsedToken.Claims.(*jwt.RegisteredClaims)
	if !claimsOk {
		return uuid.Nil, errors.New("could not retrieve claims")
	}
	user_id, errUserId := uuid.Parse(claims.Subject)
	if errUserId != nil {
		return uuid.Nil, errors.New("could not parse uuid: " + errUserId.Error())
	}
	return user_id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	str := headers.Get("Authorization")
	if str == "" {
		return "", errors.New("no or blank authorization header")
	}
	str = strings.TrimSpace(strings.ReplaceAll(str, "Bearer", ""))
	if str == "" {
		return "", errors.New("no token after bearer")
	}
	return str, nil
}
