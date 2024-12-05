package handler

import (
	"context"

	"github.com/yuyacode/AppLiftMessageApi/entity"
)

type GetMessageService interface {
	GetAllMessages(ctx context.Context, messageThreadID entity.MessageThreadID) (entity.Messages, error)
}

type RegisterOAuthService interface {
	RegisterOAuth(ctx context.Context, apiKey string) error
}
