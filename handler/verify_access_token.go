package handler

import (
	"net/http"

	"github.com/yuyacode/AppLiftMessageApi/request"
)

func VerifyAccessTokenMiddleware(vat VerifyAccessTokenService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			accessToken, err := extractAuthorizationHeader(r)
			if err != nil {
				RespondJSON(ctx, w, ErrResponse{
					Message: "invalid_token",
					Detail:  err.Error(),
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
