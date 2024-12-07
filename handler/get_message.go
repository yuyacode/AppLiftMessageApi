package handler

import (
	"encoding/json"
	"net/http"
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
	IsUnread      string           `json:"is_unread"       db:"is_unread"`
	CreatedAt     time.Time        `json:"created_at"      db:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"      db:"updated_at"`
}

func NewGetMessage(service GetMessageService, validator *validator.Validate) *GetMessage {
	return &GetMessage{
		Service:   service,
		Validator: validator,
	}
}

func (gm *GetMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		MessageThreadID entity.MessageThreadID `json:"message_thread_id" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	if err := gm.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	messages, err := gm.Service.GetAllMessages(ctx, b.MessageThreadID)
	if err != nil {
		if serviceErr, ok := err.(*ServiceError); ok {
			RespondJSON(ctx, w, &ErrResponse{
				Message: serviceErr.Error(),
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
			IsUnread:      m.IsUnread,
			CreatedAt:     m.CreatedAt,
			UpdatedAt:     m.UpdatedAt,
		})
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
