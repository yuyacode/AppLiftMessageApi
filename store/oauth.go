package store

import (
	"context"
	"database/sql"

	"github.com/yuyacode/AppLiftMessageApi/clock"
	"github.com/yuyacode/AppLiftMessageApi/entity"
)

type OAuthRepository struct {
	Clocker clock.Clocker
}

func NewOAuthRepository(clocker clock.Clocker) *OAuthRepository {
	return &OAuthRepository{
		Clocker: clocker,
	}
}

func (o *OAuthRepository) GetAPIKey(ctx context.Context, db Queryer) (string, error) {
	query := "SELECT api_key FROM message_api_keys WHERE deleted_at IS NULL LIMIT 1;"
	var apiKey string
	if err := db.GetContext(ctx, &apiKey, query); err != nil {
		return "", err
	}
	return apiKey, nil
}

func (o *OAuthRepository) GetClientID(ctx context.Context, db Queryer, param *entity.MessageAPICredential) (string, error) {
	query := "SELECT client_id FROM message_api_credentials WHERE user_id = :user_id AND deleted_at IS NULL LIMIT 1;"
	var clientID string
	if err := db.GetContext(ctx, &clientID, query, param); err != nil {
		return "", err
	}
	return clientID, nil
}

func (o *OAuthRepository) GetClientSecret(ctx context.Context, db Queryer, param *entity.MessageAPICredential) (string, error) {
	query := "SELECT client_secret FROM message_api_credentials WHERE user_id = :user_id AND deleted_at IS NULL LIMIT 1;"
	var clientSecret string
	if err := db.GetContext(ctx, &clientSecret, query, param); err != nil {
		return "", err
	}
	return clientSecret, nil
}

func (o *OAuthRepository) SearchByClientID(ctx context.Context, db Queryer, messageAPICredential *entity.MessageAPICredential) (bool, error) {
	query := "SELECT 1 FROM message_api_credentials WHERE client_id = :client_id AND deleted_at IS NULL LIMIT 1;"
	var dummy int
	if err := db.GetContext(ctx, &dummy, query, messageAPICredential); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (o *OAuthRepository) SearchByClientSecret(ctx context.Context, db Queryer, messageAPICredential *entity.MessageAPICredential) (bool, error) {
	query := "SELECT 1 FROM message_api_credentials WHERE client_secret = :client_secret AND deleted_at IS NULL LIMIT 1;"
	var dummy int
	if err := db.GetContext(ctx, &dummy, query, messageAPICredential); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (o *OAuthRepository) SaveClientIDSecret(ctx context.Context, db Execer, messageAPICredential *entity.MessageAPICredential) error {
	query := "INSERT INTO message_api_credentials (user_id, client_id, client_secret) VALUES (:user_id, :client_id, :client_secret);"
	_, err := db.NamedExecContext(ctx, query, messageAPICredential)
	if err != nil {
		return err
	}
	return nil
}

func (o *OAuthRepository) SearchByAccessToken(ctx context.Context, db Queryer, messageAPICredential *entity.MessageAPICredential) (bool, error) {
	query := "SELECT 1 FROM message_api_credentials WHERE access_token = :access_token AND deleted_at IS NULL LIMIT 1;"
	var dummy int
	if err := db.GetContext(ctx, &dummy, query, messageAPICredential); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (o *OAuthRepository) SearchByRefreshToken(ctx context.Context, db Queryer, messageAPICredential *entity.MessageAPICredential) (bool, error) {
	query := "SELECT 1 FROM message_api_credentials WHERE refresh_token = :refresh_token AND deleted_at IS NULL LIMIT 1;"
	var dummy int
	if err := db.GetContext(ctx, &dummy, query, messageAPICredential); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (o *OAuthRepository) SaveToken(ctx context.Context, db Execer, messageAPICredential *entity.MessageAPICredential) error {
	query := "UPDATE message_api_credentials SET access_token = :access_token, refresh_token = :refresh_token, expires_at = :expires_at WHERE user_id = :user_id;"
	_, err := db.NamedExecContext(ctx, query, messageAPICredential)
	if err != nil {
		return err
	}
	return nil
}

func (o *OAuthRepository) GetAccessToken(ctx context.Context, db Queryer, param *entity.MessageAPICredential) (*entity.MessageAPICredential, error) {
	query := "SELECT access_token, expires_at FROM message_api_credentials WHERE user_id = :user_id AND deleted_at IS NULL LIMIT 1;"
	var result *entity.MessageAPICredential
	if err := db.GetContext(ctx, result, query, param); err != nil {
		return nil, err
	}
	return result, nil
}

func (o *OAuthRepository) GetRefreshToken(ctx context.Context, db Queryer, param *entity.MessageAPICredential) (*entity.MessageAPICredential, error) {
	query := "SELECT refresh_token FROM message_api_credentials WHERE user_id = :user_id AND deleted_at IS NULL LIMIT 1;"
	var result *entity.MessageAPICredential
	if err := db.GetContext(ctx, result, query, param); err != nil {
		return nil, err
	}
	return result, nil
}
