package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/yuyacode/AppLiftMessageApi/credential"
	"github.com/yuyacode/AppLiftMessageApi/entity"
	"github.com/yuyacode/AppLiftMessageApi/request"
)

type RegisterOAuth struct {
	DBHandlers       map[string]*sqlx.DB
	CredentialGetter CredentialGetter
	CredentialSetter CredentialSetter
}

func (r *RegisterOAuth) RegisterOAuth(ctx context.Context, apiKey string) error {
	appKind, ok := request.GetAppKind(ctx)
	if !ok {
		return fmt.Errorf("failed to get app kind")
	}
	validAPIKey, err := r.CredentialGetter.GetAPIKey(ctx, r.DBHandlers[appKind])
	if err != nil {
		return fmt.Errorf("failed to get API Key: %w", err)
	}
	if apiKey != validAPIKey {
		return fmt.Errorf("API Key is invalid")
	}
	messageAPICredential := &entity.MessageAPICredential{}
	for i := 0; i < 5; i++ {
		clientID, err := credential.GenerateClientID()
		if err != nil {
			return fmt.Errorf("failed to generate client_id: %w", err)
		}
		messageAPICredential.ClientID = clientID
		exist, err := r.CredentialGetter.SearchByClientID(ctx, r.DBHandlers[appKind], messageAPICredential)
		if err != nil {
			return fmt.Errorf("failed to search client_id: %w", err)
		}
		if !exist {
			break
		}
		if i == 4 {
			return fmt.Errorf("failed to generate client_id 5 times")
		}
	}
	for i := 0; i < 5; i++ {
		clientSecret, err := credential.GenerateClientSecret()
		if err != nil {
			return fmt.Errorf("failed to generate client_secret: %w", err)
		}
		messageAPICredential.ClientSecret = clientSecret
		exist, err := r.CredentialGetter.SearchByClientSecret(ctx, r.DBHandlers[appKind], messageAPICredential)
		if err != nil {
			return fmt.Errorf("failed to search client_secret: %w", err)
		}
		if !exist {
			break
		}
		if i == 4 {
			return fmt.Errorf("failed to generate client_secret 5 times")
		}
	}
	userID, ok := request.GetUserID(ctx)
	if !ok {
		return fmt.Errorf("failed to get userID")
	}
	messageAPICredential.UserID = userID
	err = r.CredentialSetter.SaveClientIDSecret(ctx, r.DBHandlers[appKind], messageAPICredential)
	if err != nil {
		return fmt.Errorf("failed to insert message api client_id and client_secret: %w", err)
	}
	for i := 0; i < 5; i++ {
		access_token, err := credential.GenerateToken(appKind, userID)
		if err != nil {
			return fmt.Errorf("failed to generate access_token: %w", err)
		}
		messageAPICredential.AccessToken = access_token
		exist, err := r.CredentialGetter.SearchByAccessToken(ctx, r.DBHandlers[appKind], messageAPICredential)
		if err != nil {
			return fmt.Errorf("failed to search access_token: %w", err)
		}
		if !exist {
			break
		}
		if i == 4 {
			return fmt.Errorf("failed to generate access_token 5 times")
		}
	}
	for i := 0; i < 5; i++ {
		refresh_token, err := credential.GenerateToken(appKind, userID)
		if err != nil {
			return fmt.Errorf("failed to generate refresh_token: %w", err)
		}
		messageAPICredential.RefreshToken = refresh_token
		exist, err := r.CredentialGetter.SearchByRefreshToken(ctx, r.DBHandlers[appKind], messageAPICredential)
		if err != nil {
			return fmt.Errorf("failed to search refresh_token: %w", err)
		}
		if !exist {
			break
		}
		if i == 4 {
			return fmt.Errorf("failed to generate refresh_token 5 times")
		}
	}
	messageAPICredential.ExpiresAt = time.Now().Add(15 * time.Minute)
	err = r.CredentialSetter.SaveToken(ctx, r.DBHandlers[appKind], messageAPICredential)
	if err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}
	return nil
}
