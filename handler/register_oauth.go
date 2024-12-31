package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/yuyacode/AppLiftMessageApi/request"
)

type RegisterOAuth struct {
	Service   RegisterOAuthService
	Validator *validator.Validate
}

func NewRegisterOAuth(service RegisterOAuthService, validator *validator.Validate) *RegisterOAuth {
	return &RegisterOAuth{
		Service:   service,
		Validator: validator,
	}
}

func (ro *RegisterOAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestData struct {
		APIKey  string `                validate:"required"`
		UserID  int64  `json:"user_id"  validate:"required,numeric"`
		AppKind string `json:"app_kind" validate:"required,oneof=company student"`
	}
	apiKey, err := extractAuthorizationHeader(r)
	if err != nil {
		RespondJSON(ctx, w, ErrResponse{
			Message: "invalid_api_key",
			Detail:  err.Error(),
		}, http.StatusUnauthorized)
		return
	}
	requestData.APIKey = apiKey
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	if err := ro.Validator.Struct(requestData); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	ctx = request.SetAppKind(ctx, requestData.AppKind)
	ctx = request.SetUserID(ctx, requestData.UserID)
	err = ro.Service.RegisterOAuth(ctx, requestData.APIKey)
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
	RespondJSON(ctx, w, &SuccessResponse{
		Message: "OAuth registration was successful",
	}, http.StatusOK)
}
