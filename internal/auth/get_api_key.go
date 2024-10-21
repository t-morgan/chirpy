package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth_header := headers.Get("authorization")
	if auth_header == "" {
		return "", errors.New("no authorization header found")
	}

	return strings.TrimPrefix(auth_header, "ApiKey "), nil
}
