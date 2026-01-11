package auth

import (
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrMissingHeader
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", ErrInvalidHeader
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}
