package batch

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/yuyacode/AppLiftMessageApi/clock"
	"github.com/yuyacode/AppLiftMessageApi/config"
	"github.com/yuyacode/AppLiftMessageApi/credential"
	"github.com/yuyacode/AppLiftMessageApi/store"
)

const apiKeyLength = 32

func GenerateAPIKey(target string) (string, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	dbHandler, dbCloseFunc, err := store.New(ctx, cfg, target)
	if err != nil {
		dbCloseFunc()
		return "", err
	}
	defer dbCloseFunc()
	bytes := make([]byte, apiKeyLength)
	_, err = rand.Read(bytes)
	if err != nil {
		return "", err
	}
	apiKey := hex.EncodeToString(bytes)
	hashedAPIKey := credential.HashAPIKey(apiKey)
	clocker := clock.RealClocker{}
	query := "INSERT INTO message_api_keys (api_key, created_at) VALUES (?, ?);"
	_, err = dbHandler.ExecContext(ctx, query, hashedAPIKey, clocker.Now())
	if err != nil {
		return "", err
	}
	return apiKey, nil
}
