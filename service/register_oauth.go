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
	messageAPICredential := &entity.MessageAPICredential{}
	for i := 0; i < 5; i++ {
		clientID, err := credential.GenerateClientID()
		if err != nil {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to generate client_id",
				err.Error(),
			)
		}
		messageAPICredential.ClientID = clientID
		exist, err := ro.CredentialGetter.SearchByClientID(ctx, ro.DBHandlers[appKind], messageAPICredential)
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
	for i := 0; i < 5; i++ {
		clientSecret, err := credential.GenerateClientSecret()
		if err != nil {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to generate client_secret",
				err.Error(),
			)
		}
		messageAPICredential.ClientSecret = clientSecret
		exist, err := ro.CredentialGetter.SearchByClientSecret(ctx, ro.DBHandlers[appKind], messageAPICredential)
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
	messageAPICredential.UserID = userID
	err = ro.CredentialSetter.SaveClientIDSecret(ctx, ro.DBHandlers[appKind], messageAPICredential)
	if err != nil {
		return handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to insert message api client_id and client_secret",
			err.Error(),
		)
	}
	for i := 0; i < 5; i++ {
		access_token, err := credential.GenerateAccessToken(appKind, userID)
		if err != nil {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to generate access_token",
				err.Error(),
			)
		}
		messageAPICredential.AccessToken = access_token
		exist, err := ro.CredentialGetter.SearchByAccessToken(ctx, ro.DBHandlers[appKind], messageAPICredential)
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
	for i := 0; i < 5; i++ {
		refresh_token, err := credential.GenerateRefreshToken(appKind, userID)
		if err != nil {
			return handler.NewServiceError(
				http.StatusInternalServerError,
				"failed to generate refresh_token",
				err.Error(),
			)
		}
		messageAPICredential.RefreshToken = refresh_token
		exist, err := ro.CredentialGetter.SearchByRefreshToken(ctx, ro.DBHandlers[appKind], messageAPICredential)
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
	messageAPICredential.ExpiresAt = time.Now().Add(15 * time.Minute)
	err = ro.CredentialSetter.SaveToken(ctx, ro.DBHandlers[appKind], messageAPICredential)
	if err != nil {
		return handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to save token",
			err.Error(),
		)
	}
	return nil
}
