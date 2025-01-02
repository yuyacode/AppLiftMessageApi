package service

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/yuyacode/AppLiftMessageApi/entity"
	"github.com/yuyacode/AppLiftMessageApi/handler"
	"github.com/yuyacode/AppLiftMessageApi/request"
)

type AddMessage struct {
	DBHandlers         map[string]*sqlx.DB
	MessageAdder       MessageAdder
	MessageOwnerGetter MessageOwnerGetter
}

func NewAddMessage(dbHandlers map[string]*sqlx.DB, messageAdder MessageAdder, messageOwnerGetter MessageOwnerGetter) *AddMessage {
	return &AddMessage{
		DBHandlers:         dbHandlers,
		MessageAdder:       messageAdder,
		MessageOwnerGetter: messageOwnerGetter,
	}
}

func (am *AddMessage) AddMessage(ctx context.Context, messageThreadID entity.MessageThreadID, isFromCompany int8, isFromStudent int8, content string, isSent int8, sentAt time.Time) (*entity.Message, error) {
	appKind, ok := request.GetAppKind(ctx)
	if !ok {
		return nil, handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get app kind",
			"",
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
	if appKind == "company" {
		companyUserID, err := am.MessageOwnerGetter.GetThreadCompanyOwner(ctx, am.DBHandlers["common"], messageThreadID)
		if err != nil {
			return nil, handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to get threadCompanyOwner",
				err.Error(),
			)
		}
		if userID != companyUserID {
			return nil, handler.NewServiceError(
				http.StatusForbidden,
				"unauthorized: lack the necessary permissions to add messages",
				"",
			)
		}
	} else if appKind == "student" {
		studentUserID, err := am.MessageOwnerGetter.GetThreadStudentOwner(ctx, am.DBHandlers["common"], messageThreadID)
		if err != nil {
			return nil, handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to get threadStudentOwner",
				err.Error(),
			)
		}
		if userID != studentUserID {
			return nil, handler.NewServiceError(
				http.StatusForbidden,
				"unauthorized: lack the necessary permissions to add messages",
				"",
			)
		}
	}
	m := &entity.Message{
		MessageThreadID: messageThreadID,
		IsFromCompany:   isFromCompany,
		IsFromStudent:   isFromStudent,
		Content:         content,
		IsSent:          isSent,
		SentAt:          sentAt,
	}
	err := am.MessageAdder.AddMessage(ctx, am.DBHandlers["common"], m)
	if err != nil {
		return nil, handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to add message",
			err.Error(),
		)
	}
	return m, nil
}
