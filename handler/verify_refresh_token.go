package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/yuyacode/AppLiftMessageApi/request"
)

func VerifyRefreshTokenMiddleware(vrt VerifyRefreshTokenService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			body, err := io.ReadAll(r.Body)
			if err != nil {
				RespondJSON(ctx, w, ErrResponse{
					Message: "failed to read request body",
					Detail:  err.Error(),
				}, http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			var b map[string]interface{}
			if err := json.Unmarshal(body, &b); err != nil {
				RespondJSON(ctx, w, ErrResponse{
					Message: "invalid JSON format",
					Detail:  err.Error(),
				}, http.StatusBadRequest)
				return
			}
			refresh_token, ok := b["refresh_token"].(string)
			if !ok || refresh_token == "" {
				RespondJSON(ctx, w, ErrResponse{
					Message: "invalid_token",
					Detail:  "invalid refresh token",
				}, http.StatusUnauthorized)
				return
			}
			appKind, userID, err := vrt.VerifyRefreshToken(ctx, refresh_token)
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
