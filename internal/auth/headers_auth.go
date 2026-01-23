package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	return getValueFromAuthHeader(headers, "Bearer")
}

func GetAPIKey(headers http.Header) (string, error) {
	return getValueFromAuthHeader(headers, "ApiKey")
}

func getValueFromAuthHeader(headers http.Header, key string) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", ErrMissingHeader
	}

	keyWithWhitespace := fmt.Sprintf("%s ", key)
	if !strings.HasPrefix(authHeader, keyWithWhitespace) {
		return "", ErrInvalidHeader
	}

	return strings.TrimPrefix(authHeader, keyWithWhitespace), nil
}
