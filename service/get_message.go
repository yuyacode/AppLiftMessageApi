package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/yuyacode/AppLiftMessageApi/entity"
)

type GetMessage struct {
	DB    map[string]*sqlx.DB
	Repo  MessageGetter
	Owner MessageOwnerGetter
}

func (g *GetMessage) GetAllMessages(ctx context.Context, messageThreadID entity.MessageThreadID) (entity.Messages, error) {
	companyUserID, err := g.Owner.GetThreadCompanyOwner(ctx, g.DB["common"], messageThreadID)
	if err != nil {
		return nil, fmt.Errorf("failed to get threadCompanyOwner: %w", err)
	}
	// 認可判定
	m, err := g.Repo.GetAllMessages(ctx, g.DB["common"], messageThreadID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	return m, nil
}
