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
	DB    map[string]*sqlx.DB
	Repo  MessageGetter
	Owner MessageOwnerGetter
}

func (g *GetMessage) GetAllMessages(ctx context.Context, messageThreadID entity.MessageThreadID) (entity.Messages, error) {
	messageThread := &entity.MessageThread{ID: messageThreadID}
	companyUserID, err := g.Owner.GetThreadCompanyOwner(ctx, g.DB["common"], messageThread)
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
	m, err := g.Repo.GetAllMessages(ctx, g.DB["common"], messageThreadID)
	if err != nil {
		return nil, handler.NewServiceError(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to get message: %v", err),
		)
	}
	return m, nil
}
