package service

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/yuyacode/AppLiftMessageApi/credential"
	"github.com/yuyacode/AppLiftMessageApi/handler"
)

type VerifyAccessToken struct {
	DBHandlers       map[string]*sqlx.DB
	CredentialGetter CredentialGetter
}

func NewVerifyAccessToken(dbHandlers map[string]*sqlx.DB, credentialGetter CredentialGetter) *VerifyAccessToken {
	return &VerifyAccessToken{
		DBHandlers:       dbHandlers,
		CredentialGetter: credentialGetter,
	}
}

func (vat *VerifyAccessToken) VerifyAccessToken(ctx context.Context, accessToken string) (string, int64, error) {
	appKind, userID, err := credential.DecryptAccessToken(accessToken)
	if err != nil {
		return "", 0, err
	}
	validAccessToken, expiresAt, err := vat.CredentialGetter.GetAccessToken(ctx, vat.DBHandlers[appKind], userID)
	if err != nil {
		return "", 0, handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get access_token",
			err.Error(),
		)
	}
	if accessToken != validAccessToken {
		return "", 0, handler.NewServiceError(
			http.StatusUnauthorized,
			"invalid_token",
			"invalid access token",
		)
	}
	currentTime := time.Now()
	if currentTime.After(expiresAt) {
		return "", 0, handler.NewServiceError(
			http.StatusUnauthorized,
			"token_expired",
			"The access token has expired",
		)
	}
	return appKind, userID, nil
}
