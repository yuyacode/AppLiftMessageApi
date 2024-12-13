package handler

import (
	"fmt"
	"net/http"
	"strings"
)

func extractAuthorizationHeader(r *http.Request) (string, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}
	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		return "", fmt.Errorf("invalid Authorization header format")
	}
	authorizationHeader = strings.TrimPrefix(authorizationHeader, "Bearer ")
	if authorizationHeader == "" {
		return "", fmt.Errorf("empty Authorization header")
	}
	return authorizationHeader, nil
}
