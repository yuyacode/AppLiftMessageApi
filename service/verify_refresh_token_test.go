package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/yuyacode/AppLiftMessageApi/credential"
	"github.com/yuyacode/AppLiftMessageApi/handler"
	"github.com/yuyacode/AppLiftMessageApi/store"
)

func TestVerifyRefreshToken_VerifyRefreshToken(t *testing.T) {
	appKind := "student"
	userID := int64(1)
	refreshToken, err := credential.GenerateRefreshToken(appKind, userID)
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}
	type testCase struct {
		name string
		// 事前に credential.DecryptAccessToken で返される値を想定
		// (本テストではDecryptの処理自体はテストしない → 常に成功すると仮定)
		decryptedAppKind  string
		decryptedUserID   int64
		refreshToken      string
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
			refreshToken:     refreshToken,
			prepareGetterMock: func(m *CredentialGetterMock) {
				m.GetRefreshTokenFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "", errors.New("db error")
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get refresh_token",
		},
		{
			name:             "token mismatch => invalid_token",
			decryptedAppKind: appKind,
			decryptedUserID:  userID,
			refreshToken:     refreshToken,
			prepareGetterMock: func(m *CredentialGetterMock) {
				m.GetRefreshTokenFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "invalid-refresh-token", nil
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusUnauthorized,
			wantErrMsg:    "invalid_token",
		},
		{
			name:             "success",
			decryptedAppKind: appKind,
			decryptedUserID:  userID,
			refreshToken:     refreshToken,
			prepareGetterMock: func(m *CredentialGetterMock) {
				m.GetRefreshTokenFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return refreshToken, nil
				}
			},
			wantErr:     false,
			wantAppKind: appKind,
			wantUserID:  userID,
		},
	}
	dbHandlers := map[string]*sqlx.DB{
		"student": nil,
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			getterMock := &CredentialGetterMock{}
			if tc.prepareGetterMock != nil {
				tc.prepareGetterMock(getterMock)
			}
			svc := NewVerifyRefreshToken(dbHandlers, getterMock)
			appKind, userID, err := svc.VerifyRefreshToken(context.Background(), tc.refreshToken)
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
