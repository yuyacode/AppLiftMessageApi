package service

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/yuyacode/AppLiftMessageApi/entity"
	"github.com/yuyacode/AppLiftMessageApi/handler"
	"github.com/yuyacode/AppLiftMessageApi/request"
)

type DeleteMessage struct {
	DBHandlers         map[string]*sqlx.DB
	MessageDeleter     MessageDeleter
	MessageOwnerGetter MessageOwnerGetter
}

func NewDeleteMessage(dbHandlers map[string]*sqlx.DB, messageDeleter MessageDeleter, messageOwnerGetter MessageOwnerGetter) *DeleteMessage {
	return &DeleteMessage{
		DBHandlers:         dbHandlers,
		MessageDeleter:     messageDeleter,
		MessageOwnerGetter: messageOwnerGetter,
	}
}

func (dm *DeleteMessage) DeleteMessage(ctx context.Context, id entity.MessageID) error {
	companyUserID, err := dm.MessageOwnerGetter.GetThreadCompanyOwnerByMessageID(ctx, dm.DBHandlers["common"], id)
	if err != nil {
		return handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get threadCompanyOwner",
			err.Error(),
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
	if userID != companyUserID {
		return handler.NewServiceError(
			http.StatusForbidden,
			"unauthorized: lack the necessary permissions to delete message",
			"",
		)
	}
	err = dm.MessageDeleter.DeleteMessage(ctx, dm.DBHandlers["common"], id)
	if err != nil {
		return handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to delete message",
			err.Error(),
		)
	}
	return nil
}
