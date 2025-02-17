package handler

import (
	"context"
	"time"

	"github.com/yuyacode/AppLiftMessageApi/entity"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . GetMessageService

type VerifyAccessTokenService interface {
	VerifyAccessToken(ctx context.Context, accessToken string) (string, int64, error)
}

type VerifyRefreshTokenService interface {
	VerifyRefreshToken(ctx context.Context, refreshToken string) (string, int64, error)
}

type RegisterOAuthService interface {
	RegisterOAuth(ctx context.Context, apiKey string) error
}

type RefreshAccessTokenService interface {
	RefreshAccessToken(ctx context.Context, client_id, client_secret string) (string, string, error)
}

type GetMessageService interface {
	GetAllMessages(ctx context.Context, messageThreadID entity.MessageThreadID) (entity.Messages, error)
}

type AddMessageService interface {
	AddMessage(ctx context.Context, messageThreadID entity.MessageThreadID, isFromCompany int8, isFromStudent int8, content string, isSent int8, sentAt time.Time) (*entity.Message, error)
}

type EditMessageService interface {
	EditMessage(ctx context.Context, id entity.MessageID, content string) error
}

type DeleteMessageService interface {
	DeleteMessage(ctx context.Context, id entity.MessageID) error
}
