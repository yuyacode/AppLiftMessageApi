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

type VerifyAccessTokenService interface {
	VerifyAccessToken(ctx context.Context, accessToken string) (string, int64, error)
}
