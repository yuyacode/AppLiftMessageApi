package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yuyacode/AppLiftMessageApi/request"
)

func VerifyAccessToken(vat VerifyAccessTokenService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			accessToken, err := extractAccessToken(r)
			if err != nil {
				RespondJSON(ctx, w, ErrResponse{
					Message: "invalid_token",
					Detail:  fmt.Sprintf("not found oauth info: %v", err),
				}, http.StatusUnauthorized)
				return
			}
			appKind, userID, err := vat.VerifyAccessToken(ctx, accessToken)
			if err != nil {
				if serviceErr, ok := err.(*ServiceError); ok {
					RespondJSON(ctx, w, &ErrResponse{
						Message: serviceErr.Error(),
						Detail:  serviceErr.DetailError(),
					}, serviceErr.StatusCode)
					return
				}
				RespondJSON(ctx, w, &ErrResponse{
					Message: err.Error(),
				}, http.StatusInternalServerError)
				return
			}
			ctx = request.SetAppKind(ctx, appKind)
			ctx = request.SetUserID(ctx, userID)
			clone := r.Clone(ctx)
			next.ServeHTTP(w, clone)
		})
	}
}

func extractAccessToken(r *http.Request) (string, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}
	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		return "", fmt.Errorf("invalid Authorization header format")
	}
	accessToken := strings.TrimPrefix(authorizationHeader, "Bearer ")
	if accessToken == "" {
		return "", fmt.Errorf("empty accessToken")
	}
	return accessToken, nil
}
