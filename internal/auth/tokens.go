package auth

import (
	"errors"
	"net/http"
	"strings"
)

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

func GetAPIKey(headers http.Header) (string, error) {
	str := headers.Get("Authorization")
	if str == "" {
		return "", errors.New("no or blank authorization header")
	}
	str = strings.TrimSpace(strings.ReplaceAll(str, "ApiKey", ""))
	if str == "" {
		return "", errors.New("no token after bearer")
	}
	return str, nil
}
