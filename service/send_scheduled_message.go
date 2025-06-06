package service

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/yuyacode/AppLiftMessageApi/handler"
)

type SendScheduledMessage struct {
	DBHandlers             map[string]*sqlx.DB
	ScheduledMessageSender ScheduledMessageSender
}

func NewSendScheduledMessage(dbHandlers map[string]*sqlx.DB, ScheduledMessageSender ScheduledMessageSender) *SendScheduledMessage {
	return &SendScheduledMessage{
		DBHandlers:             dbHandlers,
		ScheduledMessageSender: ScheduledMessageSender,
	}
}

func (ssm *SendScheduledMessage) SendScheduledMessage(ctx context.Context) error {
	err := ssm.ScheduledMessageSender.SendScheduledMessage(ctx, ssm.DBHandlers["common"])
	if err != nil {
		return handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to send scheduled message",
			err.Error(),
		)
	}
	return nil
}
