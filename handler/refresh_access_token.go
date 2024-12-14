package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type RefreshAccessToken struct {
	Service   RefreshAccessTokenService
	Validator *validator.Validate
}

func NewRefreshAccessToken(service RefreshAccessTokenService, validator *validator.Validate) *RefreshAccessToken {
	return &RefreshAccessToken{
		Service:   service,
		Validator: validator,
	}
}

func (rat *RefreshAccessToken) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestData struct {
		ClientID     string `json:"client_id"     validate:"required"`
		ClientSecret string `json:"client_secret" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	if err := rat.Validator.Struct(requestData); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusUnauthorized)
		return
	}
	accessToken, err := rat.Service.RefreshAccessToken(ctx, requestData.ClientID, requestData.ClientSecret)
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
	rsp := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: accessToken,
	}
	RespondJSON(ctx, w, &rsp, http.StatusOK)
}
