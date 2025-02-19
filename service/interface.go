package service

import (
	"context"
	"database/sql"

	"github.com/yuyacode/AppLiftMessageApi/entity"
	"github.com/yuyacode/AppLiftMessageApi/store"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . MessageOwnerGetter MessageGetter

type CredentialGetter interface {
	GetAPIKey(ctx context.Context, db store.Queryer) (string, error)
	GetClientID(ctx context.Context, db store.Queryer, userID int64) (string, error)
	GetClientSecret(ctx context.Context, db store.Queryer, userID int64) (string, error)
	GetAccessToken(ctx context.Context, db store.Queryer, userID int64) (string, *sql.NullTime, error)
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

type MessageOwnerGetter interface {
	GetThreadCompanyOwner(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (int64, error)
	GetThreadStudentOwner(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (int64, error)
	GetThreadCompanyOwnerByMessageID(ctx context.Context, db store.Queryer, messageID entity.MessageID) (int64, error)
	GetThreadStudentOwnerByMessageID(ctx context.Context, db store.Queryer, messageID entity.MessageID) (int64, error)
}

type MessageGetter interface {
	GetAllMessagesForCompanyUser(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (entity.Messages, error)
	GetAllMessagesForStudentUser(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (entity.Messages, error)
}

type MessageAdder interface {
	AddMessage(ctx context.Context, db store.Execer, param *entity.Message) error
}

type MessageEditor interface {
	EditMessage(ctx context.Context, db store.Execer, param *entity.Message) error
}

type MessageDeleter interface {
	DeleteMessage(ctx context.Context, db store.Execer, id entity.MessageID) error
}
