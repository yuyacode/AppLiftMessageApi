package service

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/yuyacode/AppLiftMessageApi/entity"
	"github.com/yuyacode/AppLiftMessageApi/handler"
	"github.com/yuyacode/AppLiftMessageApi/request"
)

type GetMessage struct {
	DBHandlers         map[string]*sqlx.DB
	MessageGetter      MessageGetter
	MessageOwnerGetter MessageOwnerGetter
}

func NewGetMessage(dbHandlers map[string]*sqlx.DB, messageGetter MessageGetter, messageOwnerGetter MessageOwnerGetter) *GetMessage {
	return &GetMessage{
		DBHandlers:         dbHandlers,
		MessageGetter:      messageGetter,
		MessageOwnerGetter: messageOwnerGetter,
	}
}

func (gm *GetMessage) GetAllMessages(ctx context.Context, messageThreadID entity.MessageThreadID) (entity.Messages, error) {
	companyUserID, err := gm.MessageOwnerGetter.GetThreadCompanyOwner(ctx, gm.DBHandlers["common"], messageThreadID)
	if err != nil {
		return nil, handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get threadCompanyOwner",
			err.Error(),
		)
	}
	userID, ok := request.GetUserID(ctx)
	if !ok {
		return nil, handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get userID",
			"",
		)
	}
	if userID != companyUserID {
		return nil, handler.NewServiceError(
			http.StatusForbidden,
			"unauthorized: lack the necessary permissions to retrieve messages",
			"",
		)
	}
	m, err := gm.MessageGetter.GetAllMessagesForCompanyUser(ctx, gm.DBHandlers["common"], messageThreadID)
	if err != nil {
		return nil, handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get message",
			err.Error(),
		)
	}
	return m, nil
}
