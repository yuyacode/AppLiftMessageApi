package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

type SendScheduledMessage struct {
	Service   SendScheduledMessageService
	Validator *validator.Validate
}

func NewSendScheduledMessage(service SendScheduledMessageService, validator *validator.Validate) *SendScheduledMessage {
	return &SendScheduledMessage{
		Service:   service,
		Validator: validator,
	}
}

func (ssm *SendScheduledMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := ssm.Service.SendScheduledMessage(ctx)
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
		Message: "send scheduled message was successful",
	}, http.StatusOK)
}
