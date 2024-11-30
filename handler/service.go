package handler

import (
	"context"

	"github.com/yuyacode/AppLiftMessageApi/entity"
)

type GetMessageService interface {
	GetAllMessages(ctx context.Context, messageThreadID entity.MessageThreadID) (entity.Messages, error)
}
