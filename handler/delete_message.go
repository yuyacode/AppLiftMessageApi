package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/yuyacode/AppLiftMessageApi/entity"
)

type DeleteMessage struct {
	Service   DeleteMessageService
	Validator *validator.Validate
}

func NewDeleteMessage(service DeleteMessageService, validator *validator.Validate) *DeleteMessage {
	return &DeleteMessage{
		Service:   service,
		Validator: validator,
	}
}

func (dm *DeleteMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: "ID must be a number",
		}, http.StatusBadRequest)
		return
	}
	err = dm.Service.DeleteMessage(ctx, entity.MessageID(id))
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
		Message: "delete message was successful",
	}, http.StatusOK)
}
