package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/yuyacode/AppLiftMessageApi/entity"
)

type GetMessage struct {
	Service   GetMessageService
	Validator *validator.Validate
}

type message struct {
	ID            entity.MessageID `json:"id"              db:"id"`
	IsFromCompany int8             `json:"is_from_company" db:"is_from_company"`
	IsFromStudent int8             `json:"is_from_student" db:"is_from_student"`
	Content       string           `json:"content"         db:"content"`
	IsSent        int8             `json:"is_sent"         db:"is_sent"`
	SentAt        time.Time        `json:"sent_at"         db:"sent_at"`
}

func NewGetMessage(service GetMessageService, validator *validator.Validate) *GetMessage {
	return &GetMessage{
		Service:   service,
		Validator: validator,
	}
}

func (gm *GetMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	threadIDStr := r.URL.Query().Get("thread_id")
	if threadIDStr == "" {
		RespondJSON(ctx, w, &ErrResponse{
			Message: "missing required query parameter: thread_id",
		}, http.StatusBadRequest)
		return
	}
	threadIDInt, err := strconv.ParseInt(threadIDStr, 10, 64)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: "invalid format for query parameter: thread_id. Must be a valid integer",
		}, http.StatusInternalServerError)
		return
	}
	messages, err := gm.Service.GetAllMessages(ctx, entity.MessageThreadID(threadIDInt))
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
	rsp := []message{}
	for _, m := range messages {
		rsp = append(rsp, message{
			ID:            m.ID,
			IsFromCompany: m.IsFromCompany,
			IsFromStudent: m.IsFromStudent,
			Content:       m.Content,
			IsSent:        m.IsSent,
			SentAt:        m.SentAt,
		})
	}
	RespondJSON(ctx, w, &rsp, http.StatusOK)
}
