package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/yuyacode/AppLiftMessageApi/entity"
)

type EditMessage struct {
	Service   EditMessageService
	Validator *validator.Validate
}

func NewEditMessage(service EditMessageService, validator *validator.Validate) *EditMessage {
	return &EditMessage{
		Service:   service,
		Validator: validator,
	}
}

func (em *EditMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: "ID must be a number",
		}, http.StatusBadRequest)
		return
	}
	var requestData struct {
		Content string `json:"content" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	if err := em.Validator.Struct(requestData); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	err = em.Service.EditMessage(ctx, entity.MessageID(id), requestData.Content)
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
		Message: "edit message was successful",
	}, http.StatusOK)
}
