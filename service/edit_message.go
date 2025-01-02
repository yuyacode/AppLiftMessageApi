package service

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/yuyacode/AppLiftMessageApi/entity"
	"github.com/yuyacode/AppLiftMessageApi/handler"
	"github.com/yuyacode/AppLiftMessageApi/request"
)

type EditMessage struct {
	DBHandlers         map[string]*sqlx.DB
	MessageEditor      MessageEditor
	MessageOwnerGetter MessageOwnerGetter
}

func NewEditMessage(dbHandlers map[string]*sqlx.DB, messageEditor MessageEditor, messageOwnerGetter MessageOwnerGetter) *EditMessage {
	return &EditMessage{
		DBHandlers:         dbHandlers,
		MessageEditor:      messageEditor,
		MessageOwnerGetter: messageOwnerGetter,
	}
}

func (em *EditMessage) EditMessage(ctx context.Context, id entity.MessageID, content string) error {
	appKind, ok := request.GetAppKind(ctx)
	if !ok {
		return handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get app kind",
			"",
		)
	}
	userID, ok := request.GetUserID(ctx)
	if !ok {
		return handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get userID",
			"",
		)
	}
	if appKind == "company" {
		companyUserID, err := em.MessageOwnerGetter.GetThreadCompanyOwnerByMessageID(ctx, em.DBHandlers["common"], id)
		if err != nil {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to get threadCompanyOwner",
				err.Error(),
			)
		}
		if userID != companyUserID {
			return handler.NewServiceError(
				http.StatusForbidden,
				"unauthorized: lack the necessary permissions to edit message",
				"",
			)
		}
	} else if appKind == "student" {
		studentUserID, err := em.MessageOwnerGetter.GetThreadStudentOwnerByMessageID(ctx, em.DBHandlers["common"], id)
		if err != nil {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to get threadStudentOwner",
				err.Error(),
			)
		}
		if userID != studentUserID {
			return handler.NewServiceError(
				http.StatusForbidden,
				"unauthorized: lack the necessary permissions to edit message",
				"",
			)
		}
	}
	m := &entity.Message{
		ID:      id,
		Content: content,
	}
	err := em.MessageEditor.EditMessage(ctx, em.DBHandlers["common"], m)
	if err != nil {
		return handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to edit message",
			err.Error(),
		)
	}
	return nil
}
