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

func (ro *RegisterOAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		APIKey  string `json:"api_key"  validate:"required"`
		UserID  int64  `json:"user_id"  validate:"required"`
		AppKind string `json:"app_kind" validate:"required,oneof=company student"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	if err := ro.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	ctx = request.SetAppKind(ctx, b.AppKind)
	ctx = request.SetUserID(ctx, b.UserID)
	err := ro.Service.RegisterOAuth(ctx, b.APIKey)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	RespondJSON(ctx, w, &SuccessResponse{
		Message: "OAuth registration was successful",
	}, http.StatusOK)
}
