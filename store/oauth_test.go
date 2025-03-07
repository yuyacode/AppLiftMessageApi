package store

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yuyacode/AppLiftMessageApi/clock"
)

func TestOAuthRepository_GetAPIKey(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	or := NewOAuthRepository(clock.FixedClocker{})
	tests := map[string]struct {
		mockSetup  func()
		wantErr    bool
		wantAPIKey string
	}{
		"DB error": {
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT api_key FROM message_api_keys WHERE deleted_at IS NULL LIMIT 1;$`).
					WillReturnError(assertAnError())
			},
			wantErr:    true,
			wantAPIKey: "",
		},
		"Success": {
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT api_key FROM message_api_keys WHERE deleted_at IS NULL LIMIT 1;$`).
					WillReturnRows(sqlmock.NewRows([]string{"api_key"}).AddRow("SECRET_API_KEY"))
			},
			wantErr:    false,
			wantAPIKey: "SECRET_API_KEY",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			got, err := or.GetAPIKey(context.Background(), sqlxDB)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantAPIKey, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOAuthRepository_GetClientID(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	or := NewOAuthRepository(clock.FixedClocker{})
	tests := map[string]struct {
		userID       int64
		mockSetup    func()
		wantErr      bool
		wantClientID string
	}{
		"DB error": {
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT client_id FROM message_api_credentials WHERE user_id = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs(int64(1)).
					WillReturnError(assertAnError())
			},
			wantErr:      true,
			wantClientID: "",
		},
		"Success": {
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT client_id FROM message_api_credentials WHERE user_id = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"client_id"}).AddRow("SECRET_CLIENT_ID"))
			},
			wantErr:      false,
			wantClientID: "SECRET_CLIENT_ID",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			got, err := or.GetClientID(context.Background(), sqlxDB, tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantClientID, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOAuthRepository_GetClientSecret(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	or := NewOAuthRepository(clock.FixedClocker{})
	tests := map[string]struct {
		userID           int64
		mockSetup        func()
		wantErr          bool
		wantClientSecret string
	}{
		"DB error": {
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT client_secret FROM message_api_credentials WHERE user_id = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs(int64(1)).
					WillReturnError(assertAnError())
			},
			wantErr:          true,
			wantClientSecret: "",
		},
		"Success": {
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT client_secret FROM message_api_credentials WHERE user_id = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"client_secret"}).AddRow("SECRET_CLIENT_SECRET"))
			},
			wantErr:          false,
			wantClientSecret: "SECRET_CLIENT_SECRET",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			got, err := or.GetClientSecret(context.Background(), sqlxDB, tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantClientSecret, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOAuthRepository_GetAccessToken(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	or := NewOAuthRepository(clock.FixedClocker{})
	tests := map[string]struct {
		userID          int64
		mockSetup       func()
		wantErr         bool
		wantAccessToken string
		wantExpiresAt   *sql.NullTime
	}{
		"DB error": {
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT access_token, expires_at FROM message_api_credentials WHERE user_id = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs(int64(1)).
					WillReturnError(assertAnError())
			},
			wantErr:         true,
			wantAccessToken: "",
			wantExpiresAt:   &sql.NullTime{},
		},
		"Success": {
			userID: 2,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT access_token, expires_at FROM message_api_credentials WHERE user_id = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs(int64(2)).
					WillReturnRows(
						sqlmock.NewRows([]string{"access_token", "expires_at"}).
							AddRow("ACCESS_TOKEN_VALID", time.Date(2025, 1, 1, 12, 34, 56, 0, time.UTC)),
					)
			},
			wantErr:         false,
			wantAccessToken: "ACCESS_TOKEN_VALID",
			wantExpiresAt: &sql.NullTime{
				Time:  time.Date(2025, 1, 1, 12, 34, 56, 0, time.UTC),
				Valid: true,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			gotToken, gotExpiresAt, err := or.GetAccessToken(context.Background(), sqlxDB, tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Empty(t, gotToken)
				assert.NotNil(t, gotExpiresAt)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantAccessToken, gotToken)
				require.NotNil(t, gotExpiresAt)
				assert.Equal(t, tc.wantExpiresAt.Valid, gotExpiresAt.Valid)
				assert.Equal(t, tc.wantExpiresAt.Time, gotExpiresAt.Time)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
