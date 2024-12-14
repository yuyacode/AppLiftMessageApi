package service

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/yuyacode/AppLiftMessageApi/credential"
	"github.com/yuyacode/AppLiftMessageApi/entity"
	"github.com/yuyacode/AppLiftMessageApi/handler"
)

type VerifyRefreshToken struct {
	DBHandlers       map[string]*sqlx.DB
	CredentialGetter CredentialGetter
}

func NewVerifyRefreshToken(dbHandlers map[string]*sqlx.DB, credentialGetter CredentialGetter) *VerifyRefreshToken {
	return &VerifyRefreshToken{
		DBHandlers:       dbHandlers,
		CredentialGetter: credentialGetter,
	}
}

func (vrt *VerifyRefreshToken) VerifyRefreshToken(ctx context.Context, refreshToken string) (string, int64, error) {
	appKind, userID, err := credential.DecryptRefreshToken(refreshToken)
	if err != nil {
		return "", 0, err
	}
	param := &entity.MessageAPICredential{
		UserID: userID,
	}
	validRefreshToken, err := vrt.CredentialGetter.GetRefreshToken(ctx, vrt.DBHandlers[appKind], param)
	if err != nil {
		return "", 0, handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get refresh_token",
			err.Error(),
		)
	}
	if refreshToken != validRefreshToken.RefreshToken {
		return "", 0, handler.NewServiceError(
			http.StatusUnauthorized,
			"invalid_token",
			"invalid refresh token",
		)
	}
	return appKind, userID, nil
}
