package service

import (
	"context"
	"fmt"
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

func (g *GetMessage) GetAllMessages(ctx context.Context, messageThreadID entity.MessageThreadID) (entity.Messages, error) {
	messageThread := &entity.MessageThread{ID: messageThreadID}
	companyUserID, err := g.MessageOwnerGetter.GetThreadCompanyOwner(ctx, g.DBHandlers["common"], messageThread)
	if err != nil {
		return nil, handler.NewServiceError(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to get threadCompanyOwner: %v", err),
		)
	}
	userID, ok := request.GetUserID(ctx)
	if !ok {
		return nil, handler.NewServiceError(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to get userID"),
		)
	}
	if userID != companyUserID {
		return nil, handler.NewServiceError(
			http.StatusForbidden,
			fmt.Sprintf("unauthorized: lack the necessary permissions to retrieve messages"),
		)
	}
	m, err := g.MessageGetter.GetAllMessages(ctx, g.DBHandlers["common"], messageThreadID)
	if err != nil {
		return nil, handler.NewServiceError(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to get message: %v", err),
		)
	}
	return m, nil
}
