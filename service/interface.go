package service

import (
	"context"
	"time"

	"github.com/yuyacode/AppLiftMessageApi/entity"
	"github.com/yuyacode/AppLiftMessageApi/store"
)

type CredentialGetter interface {
	GetAPIKey(ctx context.Context, db store.Queryer) (string, error)
	GetClientID(ctx context.Context, db store.Queryer, userID int64) (string, error)
	GetClientSecret(ctx context.Context, db store.Queryer, userID int64) (string, error)
	GetAccessToken(ctx context.Context, db store.Queryer, userID int64) (string, time.Time, error)
	GetRefreshToken(ctx context.Context, db store.Queryer, userID int64) (string, error)
	SearchByClientID(ctx context.Context, db store.Queryer, clientID string) (bool, error)
	SearchByClientSecret(ctx context.Context, db store.Queryer, clientSecret string) (bool, error)
	SearchByAccessToken(ctx context.Context, db store.Queryer, accessToken string) (bool, error)
	SearchByRefreshToken(ctx context.Context, db store.Queryer, refreshToken string) (bool, error)
}

type CredentialSetter interface {
	SaveClientIDSecret(ctx context.Context, db store.Execer, param *entity.MessageAPICredential) error
	SaveToken(ctx context.Context, db store.Execer, param *entity.MessageAPICredential) error
}

type MessageGetter interface {
	GetAllMessages(ctx context.Context, db store.Queryer, threadID entity.MessageThreadID) (entity.Messages, error)
}

type MessageOwnerGetter interface {
	GetThreadCompanyOwner(ctx context.Context, db store.Queryer, param *entity.MessageThread) (int64, error)
}
