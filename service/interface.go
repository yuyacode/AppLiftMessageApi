package service

import (
	"context"

	"github.com/yuyacode/AppLiftMessageApi/entity"
	"github.com/yuyacode/AppLiftMessageApi/store"
)

type MessageGetter interface {
	GetAllMessages(ctx context.Context, db store.Queryer, threadID entity.MessageThreadID) (entity.Messages, error)
}

type MessageOwnerGetter interface {
	GetThreadCompanyOwner(ctx context.Context, db store.Queryer, messageThread *entity.MessageThread) (int64, error)
}

type CredentialGetter interface {
	GetAPIKey(ctx context.Context, db store.Queryer) (string, error)
	SearchByClientID(ctx context.Context, db store.Queryer, messageAPICredential *entity.MessageAPICredential) (bool, error)
	SearchByClientSecret(ctx context.Context, db store.Queryer, messageAPICredential *entity.MessageAPICredential) (bool, error)
	SearchByAccessToken(ctx context.Context, db store.Queryer, messageAPICredential *entity.MessageAPICredential) (bool, error)
	SearchByRefreshToken(ctx context.Context, db store.Queryer, messageAPICredential *entity.MessageAPICredential) (bool, error)
	GetAccessToken(ctx context.Context, db store.Queryer, param *entity.MessageAPICredential) (*entity.MessageAPICredential, error)
}

type CredentialSetter interface {
	SaveClientIDSecret(ctx context.Context, db store.Execer, messageAPICredential *entity.MessageAPICredential) error
	SaveToken(ctx context.Context, db store.Execer, messageAPICredential *entity.MessageAPICredential) error
}
