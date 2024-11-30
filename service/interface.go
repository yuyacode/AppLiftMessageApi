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
	GetThreadCompanyOwner(ctx context.Context, db store.Queryer, threadID entity.MessageThreadID) (int64, error)
}
