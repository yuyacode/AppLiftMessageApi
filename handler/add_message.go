package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/yuyacode/AppLiftMessageApi/entity"
)

type AddMessage struct {
	Service   AddMessageService
	Validator *validator.Validate
}

func NewAddMessage(service AddMessageService, validator *validator.Validate) *AddMessage {
	return &AddMessage{
		Service:   service,
		Validator: validator,
	}
}

func (am *AddMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestData struct {
		MessageThreadID entity.MessageThreadID `json:"message_thread_id" validate:"required,numeric"`
		IsFromCompany   int8                   `json:"is_from_company"   validate:"oneof=0 1"`
		IsFromStudent   int8                   `json:"is_from_student"   validate:"oneof=0 1"`
		Content         string                 `json:"content"           validate:"required"`
		IsSent          int8                   `json:"is_sent"           validate:"oneof=0 1"`
		SentAt          time.Time              `json:"sent_at"           validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	if err := am.Validator.Struct(requestData); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	message, err := am.Service.AddMessage(ctx, requestData.MessageThreadID, requestData.IsFromCompany, requestData.IsFromStudent, requestData.Content, requestData.IsSent, requestData.SentAt)
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
		ID entity.MessageID `json:"id"`
	}{ID: message.ID}
	RespondJSON(ctx, w, &rsp, http.StatusOK)
}
