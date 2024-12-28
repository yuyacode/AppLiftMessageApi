package handler

import (
	"database/sql"
	"net/http"
	"strconv"

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
	CreatedAt     *sql.NullTime    `json:"created_at"      db:"created_at"`
	UpdatedAt     *sql.NullTime    `json:"updated_at"      db:"updated_at"`
}

func NewGetMessage(service GetMessageService, validator *validator.Validate) *GetMessage {
	return &GetMessage{
		Service:   service,
		Validator: validator,
	}
}

func (gm *GetMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var qp struct {
		MessageThreadID entity.MessageThreadID `validate:"required"`
	}
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
	qp.MessageThreadID = entity.MessageThreadID(threadIDInt)
	if err := gm.Validator.Struct(qp); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	messages, err := gm.Service.GetAllMessages(ctx, qp.MessageThreadID)
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
			IsUnread:      m.IsUnread,
			CreatedAt:     m.CreatedAt,
			UpdatedAt:     m.UpdatedAt,
		})
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
