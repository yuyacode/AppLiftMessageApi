package handler

import (
	"context"

	"github.com/yuyacode/AppLiftMessageApi/entity"
)

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
