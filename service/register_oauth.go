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

type RegisterOAuth struct {
	DBHandlers       map[string]*sqlx.DB
	CredentialGetter CredentialGetter
	CredentialSetter CredentialSetter
}

func NewRegisterOAuth(dbHandlers map[string]*sqlx.DB, credentialGetter CredentialGetter, credentialSetter CredentialSetter) *RegisterOAuth {
	return &RegisterOAuth{
		DBHandlers:       dbHandlers,
		CredentialGetter: credentialGetter,
		CredentialSetter: credentialSetter,
	}
}

func (ro *RegisterOAuth) RegisterOAuth(ctx context.Context, apiKey string) error {
	appKind, ok := request.GetAppKind(ctx)
	if !ok {
		return handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get app kind",
			"",
		)
	}
	validAPIKey, err := ro.CredentialGetter.GetAPIKey(ctx, ro.DBHandlers[appKind])
	if err != nil {
		return handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get API Key",
			err.Error(),
		)
	}
	hashedAPIKey := credential.HashAPIKey(apiKey)
	if hashedAPIKey != validAPIKey {
		return handler.NewServiceError(
			http.StatusUnauthorized,
			"API Key is invalid",
			"",
		)
	}
	var clientID string
	for i := 0; i < 5; i++ {
		var err error
		clientID, err = credential.GenerateClientID()
		if err != nil {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to generate client_id",
				err.Error(),
			)
		}
		exist, err := ro.CredentialGetter.SearchByClientID(ctx, ro.DBHandlers[appKind], clientID)
		if err != nil {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to search client_id",
				err.Error(),
			)
		}
		if !exist {
			break
		}
		if i == 4 {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to generate client_id 5 times",
				"",
			)
		}
	}
	var clientSecret string
	for i := 0; i < 5; i++ {
		var err error
		clientSecret, err = credential.GenerateClientSecret()
		if err != nil {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to generate client_secret",
				err.Error(),
			)
		}
		exist, err := ro.CredentialGetter.SearchByClientSecret(ctx, ro.DBHandlers[appKind], clientSecret)
		if err != nil {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to search client_secret",
				err.Error(),
			)
		}
		if !exist {
			break
		}
		if i == 4 {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to generate client_secret 5 times",
				"",
			)
		}
	}
	userID, ok := request.GetUserID(ctx)
	if !ok {
		return handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get userID",
			"",
		)
	}
	param := &entity.MessageAPICredential{
		UserID:       userID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
	err = ro.CredentialSetter.SaveClientIDSecret(ctx, ro.DBHandlers[appKind], param)
	if err != nil {
		return handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to insert message api client_id and client_secret",
			err.Error(),
		)
	}
	var accessToken string
	for i := 0; i < 5; i++ {
		var err error
		accessToken, err = credential.GenerateAccessToken(appKind, userID)
		if err != nil {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to generate access_token",
				err.Error(),
			)
		}
		exist, err := ro.CredentialGetter.SearchByAccessToken(ctx, ro.DBHandlers[appKind], accessToken)
		if err != nil {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to search access_token",
				err.Error(),
			)
		}
		if !exist {
			break
		}
		if i == 4 {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to generate access_token 5 times",
				"",
			)
		}
	}
	var refreshToken string
	for i := 0; i < 5; i++ {
		var err error
		refreshToken, err = credential.GenerateRefreshToken(appKind, userID)
		if err != nil {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to generate refresh_token",
				err.Error(),
			)
		}
		exist, err := ro.CredentialGetter.SearchByRefreshToken(ctx, ro.DBHandlers[appKind], refreshToken)
		if err != nil {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to search refresh_token",
				err.Error(),
			)
		}
		if !exist {
			break
		}
		if i == 4 {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to generate refresh_token 5 times",
				"",
			)
		}
	}
	param.AccessToken = accessToken
	param.RefreshToken = refreshToken
	param.ExpiresAt = time.Now().Add(15 * time.Minute)
	err = ro.CredentialSetter.SaveToken(ctx, ro.DBHandlers[appKind], param)
	if err != nil {
		return handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to save token",
			err.Error(),
		)
	}
	return nil
}
