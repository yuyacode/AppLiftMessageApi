package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/yuyacode/AppLiftMessageApi/entity"
	"github.com/yuyacode/AppLiftMessageApi/handler"
)

type GetMessage struct {
	DB    map[string]*sqlx.DB
	Repo  MessageGetter
	Owner MessageOwnerGetter
}

func (g *GetMessage) GetAllMessages(ctx context.Context, messageThreadID entity.MessageThreadID) (entity.Messages, error) {
	companyUserID, err := g.Owner.GetThreadCompanyOwner(ctx, g.DB["common"], messageThreadID)
	if err != nil {
		return nil, handler.NewServiceError(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to get threadCompanyOwner: %v", err),
		)
	}
	// 認可判定
	m, err := g.Repo.GetAllMessages(ctx, g.DB["common"], messageThreadID)
	if err != nil {
		return nil, handler.NewServiceError(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to get message: %v", err),
		)
	}
	return m, nil
}
