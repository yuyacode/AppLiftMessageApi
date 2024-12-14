package service

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/yuyacode/AppLiftMessageApi/credential"
	"github.com/yuyacode/AppLiftMessageApi/entity"
	"github.com/yuyacode/AppLiftMessageApi/handler"
	"github.com/yuyacode/AppLiftMessageApi/request"
)

type RefreshAccessToken struct {
	DBHandlers       map[string]*sqlx.DB
	CredentialGetter CredentialGetter
	CredentialSetter CredentialSetter
}

func NewRefreshAccessToken(dbHandlers map[string]*sqlx.DB, credentialGetter CredentialGetter, credentialSetter CredentialSetter) *RefreshAccessToken {
	return &RefreshAccessToken{
		DBHandlers:       dbHandlers,
		CredentialGetter: credentialGetter,
		CredentialSetter: credentialSetter,
	}
}

func (rat *RefreshAccessToken) RefreshAccessToken(ctx context.Context, client_id, client_secret string) (string, error) {
	appKind, ok := request.GetAppKind(ctx)
	if !ok {
		return "", handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get app kind",
			"",
		)
	}
	userID, ok := request.GetUserID(ctx)
	if !ok {
		return "", handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get user_id",
			"",
		)
	}
	param := &entity.MessageAPICredential{
		UserID: userID,
	}
	validClientID, err := rat.CredentialGetter.GetClientID(ctx, rat.DBHandlers[appKind], param)
	if err != nil {
		return "", handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get client_id",
			err.Error(),
		)
	}
	if client_id != validClientID {
		return "", handler.NewServiceError(
			http.StatusUnauthorized,
			"client_id is invalid",
			"",
		)
	}
	validClientSecret, err := rat.CredentialGetter.GetClientSecret(ctx, rat.DBHandlers[appKind], param)
	if err != nil {
		return "", handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get client_secret",
			err.Error(),
		)
	}
	if client_secret != validClientSecret {
		return "", handler.NewServiceError(
			http.StatusUnauthorized,
			"client_secret is invalid",
			"",
		)
	}
	for i := 0; i < 5; i++ {
		access_token, err := credential.GenerateAccessToken(appKind, userID)
		if err != nil {
			return "", handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to generate access_token",
				err.Error(),
			)
		}
		param.AccessToken = access_token
		exist, err := rat.CredentialGetter.SearchByAccessToken(ctx, rat.DBHandlers[appKind], param)
		if err != nil {
			return "", handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to search access_token",
				err.Error(),
			)
		}
		if !exist {
			break
		}
		if i == 4 {
			return "", handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to generate access_token 5 times",
				"",
			)
		}
	}
	for i := 0; i < 5; i++ {
		refresh_token, err := credential.GenerateRefreshToken(appKind, userID)
		if err != nil {
			return "", handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to generate refresh_token",
				err.Error(),
			)
		}
		param.RefreshToken = refresh_token
		exist, err := rat.CredentialGetter.SearchByRefreshToken(ctx, rat.DBHandlers[appKind], param)
		if err != nil {
			return "", handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to search refresh_token",
				err.Error(),
			)
		}
		if !exist {
			break
		}
		if i == 4 {
			return "", handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to generate refresh_token 5 times",
				"",
			)
		}
	}
	param.ExpiresAt = time.Now().Add(15 * time.Minute)
	err = rat.CredentialSetter.SaveToken(ctx, rat.DBHandlers[appKind], param)
	if err != nil {
		return "", handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to save token",
			err.Error(),
		)
	}
	return param.AccessToken, nil
}
