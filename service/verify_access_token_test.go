package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/yuyacode/AppLiftMessageApi/credential"
	"github.com/yuyacode/AppLiftMessageApi/handler"
	"github.com/yuyacode/AppLiftMessageApi/store"
)

func TestVerifyAccessToken_VerifyAccessToken(t *testing.T) {
	appKind := "company"
	userID := int64(1)
	accessToken, err := credential.GenerateAccessToken(appKind, userID)
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}
	type testCase struct {
		name string
		// 事前に credential.DecryptAccessToken で返される値を想定
		// (本テストではDecryptの処理自体はテストしない → 常に成功すると仮定)
		decryptedAppKind  string
		decryptedUserID   int64
		accessToken       string
		prepareGetterMock func(*CredentialGetterMock)
		wantErr           bool
		wantErrStatus     int
		wantErrMsg        string
		wantAppKind       string
		wantUserID        int64
	}
	tests := []testCase{
		{
			name:             "DB returns error",
			decryptedAppKind: appKind,
			decryptedUserID:  userID,
			accessToken:      accessToken,
			prepareGetterMock: func(m *CredentialGetterMock) {
				m.GetAccessTokenFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, *sql.NullTime, error) {
					return "", nil, errors.New("db error")
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get access_token",
		},
		{
			name:             "expiresAt is nil or invalid",
			decryptedAppKind: appKind,
			decryptedUserID:  userID,
			accessToken:      accessToken,
			prepareGetterMock: func(m *CredentialGetterMock) {
				m.GetAccessTokenFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, *sql.NullTime, error) {
					invalidTime := &sql.NullTime{Valid: false}
					return accessToken, invalidTime, nil
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "the expiresAt is not set.",
		},
		{
			name:             "token mismatch => invalid_token",
			decryptedAppKind: appKind,
			decryptedUserID:  userID,
			accessToken:      accessToken,
			prepareGetterMock: func(m *CredentialGetterMock) {
				m.GetAccessTokenFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, *sql.NullTime, error) {
					return "invalid-access-token", &sql.NullTime{
						Time:  time.Now().Add(15 * time.Minute),
						Valid: true,
					}, nil
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusUnauthorized,
			wantErrMsg:    "invalid_token",
		},
		{
			name:             "token expired => token_expired",
			decryptedAppKind: appKind,
			decryptedUserID:  userID,
			accessToken:      accessToken,
			prepareGetterMock: func(m *CredentialGetterMock) {
				m.GetAccessTokenFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, *sql.NullTime, error) {
					return accessToken, &sql.NullTime{
						Time:  time.Now().Add(-1 * time.Minute),
						Valid: true,
					}, nil
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusUnauthorized,
			wantErrMsg:    "token_expired",
		},
		{
			name:             "success",
			decryptedAppKind: appKind,
			decryptedUserID:  userID,
			accessToken:      accessToken,
			prepareGetterMock: func(m *CredentialGetterMock) {
				m.GetAccessTokenFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, *sql.NullTime, error) {
					return accessToken, &sql.NullTime{
						Time:  time.Now().Add(15 * time.Minute),
						Valid: true,
					}, nil
				}
			},
			wantErr:     false,
			wantAppKind: appKind,
			wantUserID:  userID,
		},
	}
	dbHandlers := map[string]*sqlx.DB{
		"company": nil,
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			getterMock := &CredentialGetterMock{}
			if tc.prepareGetterMock != nil {
				tc.prepareGetterMock(getterMock)
			}
			svc := NewVerifyAccessToken(dbHandlers, getterMock)
			appKind, userID, err := svc.VerifyAccessToken(context.Background(), tc.accessToken)
			if tc.wantErr {
				assert.Error(t, err, "error is expected but got nil")
				se, ok := err.(*handler.ServiceError)
				if assert.True(t, ok, "error should be *handler.ServiceError") {
					assert.Equal(t, tc.wantErrStatus, se.StatusCode)
					assert.Contains(t, se.Message, tc.wantErrMsg)
				}
				assert.Empty(t, appKind)
				assert.Equal(t, int64(0), userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantAppKind, appKind)
				assert.Equal(t, tc.wantUserID, userID)
			}
		})
	}
}
