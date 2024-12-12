package service

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/yuyacode/AppLiftMessageApi/credential"
	"github.com/yuyacode/AppLiftMessageApi/entity"
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
	param := &entity.MessageAPICredential{
		UserID: userID,
	}
	validAccessToken, err := vat.CredentialGetter.GetAccessToken(ctx, vat.DBHandlers[appKind], param)
	if err != nil {
		return "", 0, handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get access_token",
			err.Error(),
		)
	}
	if accessToken != validAccessToken.AccessToken {
		return "", 0, handler.NewServiceError(
			http.StatusUnauthorized,
			"invalid_token",
			"invalid access token",
		)
	}
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return "", 0, handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get timezone",
			err.Error(),
		)
	}
	expiresAtJST := validAccessToken.ExpiresAt.In(loc) // 本来はこの変換を不要にしたい。DBから取ってきた時点でJSTになっているか後ほど確認
	currentTimeJST := time.Now().In(loc)
	if currentTimeJST.After(expiresAtJST) {
		return "", 0, handler.NewServiceError(
			http.StatusUnauthorized,
			"token_expired",
			"The access token has expired",
		)
	}
	return appKind, userID, nil
}
