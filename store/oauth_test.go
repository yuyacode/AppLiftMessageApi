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
	"github.com/yuyacode/AppLiftMessageApi/entity"
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

func TestOAuthRepository_GetRefreshToken(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	or := NewOAuthRepository(clock.FixedClocker{})
	tests := map[string]struct {
		userID           int64
		mockSetup        func()
		wantErr          bool
		wantRefreshToken string
	}{
		"DB error": {
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT refresh_token FROM message_api_credentials WHERE user_id = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs(int64(1)).
					WillReturnError(assertAnError())
			},
			wantErr:          true,
			wantRefreshToken: "",
		},
		"Success": {
			userID: 2,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT refresh_token FROM message_api_credentials WHERE user_id = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs(int64(2)).
					WillReturnRows(
						sqlmock.NewRows([]string{"refresh_token"}).AddRow("REFRESH_TOKEN_VALID"),
					)
			},
			wantErr:          false,
			wantRefreshToken: "REFRESH_TOKEN_VALID",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			got, err := or.GetRefreshToken(context.Background(), sqlxDB, tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantRefreshToken, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOAuthRepository_SearchByClientID(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	or := NewOAuthRepository(clock.FixedClocker{})
	tests := map[string]struct {
		clientID  string
		mockSetup func()
		wantExist bool
		wantErr   bool
	}{
		"DB error": {
			clientID: "client_id_123",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT 1 FROM message_api_credentials WHERE client_id = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs("client_id_123").
					WillReturnError(assertAnError())
			},
			wantExist: false,
			wantErr:   true,
		},
		"No rows": {
			clientID: "non_exist_client_id",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT 1 FROM message_api_credentials WHERE client_id = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs("non_exist_client_id").
					WillReturnRows(sqlmock.NewRows([]string{"1"}))
			},
			wantExist: false,
			wantErr:   false,
		},
		"Found row": {
			clientID: "exist_client_id",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT 1 FROM message_api_credentials WHERE client_id = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs("exist_client_id").
					WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))
			},
			wantExist: true,
			wantErr:   false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			gotExist, err := or.SearchByClientID(context.Background(), sqlxDB, tc.clientID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantExist, gotExist)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOAuthRepository_SearchByClientSecret(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	or := NewOAuthRepository(clock.FixedClocker{})
	tests := map[string]struct {
		clientSecret string
		mockSetup    func()
		wantExist    bool
		wantErr      bool
	}{
		"DB error": {
			clientSecret: "client_secret_123",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT 1 FROM message_api_credentials WHERE client_secret = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs("client_secret_123").
					WillReturnError(assertAnError())
			},
			wantExist: false,
			wantErr:   true,
		},
		"No rows": {
			clientSecret: "non_exist_client_secret",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT 1 FROM message_api_credentials WHERE client_secret = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs("non_exist_client_secret").
					WillReturnRows(sqlmock.NewRows([]string{"1"}))
			},
			wantExist: false,
			wantErr:   false,
		},
		"Found row": {
			clientSecret: "exist_client_secret",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT 1 FROM message_api_credentials WHERE client_secret = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs("exist_client_secret").
					WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))
			},
			wantExist: true,
			wantErr:   false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			gotExist, err := or.SearchByClientSecret(context.Background(), sqlxDB, tc.clientSecret)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantExist, gotExist)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOAuthRepository_SearchByAccessToken(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	or := NewOAuthRepository(clock.FixedClocker{})
	tests := map[string]struct {
		accessToken string
		mockSetup   func()
		wantFound   bool
		wantErr     bool
	}{
		"DB error": {
			accessToken: "some_access_token",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT 1 FROM message_api_credentials WHERE access_token = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs("some_access_token").
					WillReturnError(assertAnError())
			},
			wantFound: false,
			wantErr:   true,
		},
		"No rows": {
			accessToken: "nonexistent_token",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT 1 FROM message_api_credentials WHERE access_token = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs("nonexistent_token").
					WillReturnRows(sqlmock.NewRows([]string{"1"}))
			},
			wantFound: false,
			wantErr:   false,
		},
		"Found row": {
			accessToken: "existing_token",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT 1 FROM message_api_credentials WHERE access_token = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs("existing_token").
					WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))
			},
			wantFound: true,
			wantErr:   false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			gotFound, err := or.SearchByAccessToken(context.Background(), sqlxDB, tc.accessToken)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantFound, gotFound)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOAuthRepository_SearchByRefreshToken(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	or := NewOAuthRepository(clock.FixedClocker{})
	tests := map[string]struct {
		refreshToken string
		mockSetup    func()
		wantFound    bool
		wantErr      bool
	}{
		"DB error": {
			refreshToken: "some_refresh_token",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT 1 FROM message_api_credentials WHERE refresh_token = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs("some_refresh_token").
					WillReturnError(assertAnError())
			},
			wantFound: false,
			wantErr:   true,
		},
		"No rows": {
			refreshToken: "nonexistent_token",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT 1 FROM message_api_credentials WHERE refresh_token = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs("nonexistent_token").
					WillReturnRows(sqlmock.NewRows([]string{"1"}))
			},
			wantFound: false,
			wantErr:   false,
		},
		"Found row": {
			refreshToken: "existing_token",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT 1 FROM message_api_credentials WHERE refresh_token = \? AND deleted_at IS NULL LIMIT 1;$`).
					WithArgs("existing_token").
					WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))
			},
			wantFound: true,
			wantErr:   false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			gotFound, err := or.SearchByRefreshToken(context.Background(), sqlxDB, tc.refreshToken)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantFound, gotFound)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOAuthRepository_SaveClientIDSecret(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	or := NewOAuthRepository(clock.FixedClocker{})
	tests := map[string]struct {
		inputParam *entity.MessageAPICredential
		mockSetup  func(*entity.MessageAPICredential)
		wantErr    bool
	}{
		"DB error": {
			inputParam: &entity.MessageAPICredential{
				UserID:       1,
				ClientID:     "CLIENT_ID",
				ClientSecret: "CLIENT_SECRET",
				CreatedAt:    clock.FixedClocker{}.Now(),
			},
			mockSetup: func(param *entity.MessageAPICredential) {
				mock.ExpectExec(`^INSERT INTO message_api_credentials \(user_id, client_id, client_secret, created_at\) VALUES \(\?, \?, \?, \?\);$`).
					WithArgs(
						param.UserID,
						param.ClientID,
						param.ClientSecret,
						param.CreatedAt,
					).
					WillReturnError(assertAnError())
			},
			wantErr: true,
		},
		"Success": {
			inputParam: &entity.MessageAPICredential{
				UserID:       1,
				ClientID:     "CLIENT_ID",
				ClientSecret: "CLIENT_SECRET",
				CreatedAt:    clock.FixedClocker{}.Now(),
			},
			mockSetup: func(param *entity.MessageAPICredential) {
				mock.ExpectExec(`^INSERT INTO message_api_credentials \(user_id, client_id, client_secret, created_at\) VALUES \(\?, \?, \?, \?\);$`).
					WithArgs(
						param.UserID,
						param.ClientID,
						param.ClientSecret,
						param.CreatedAt,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup(tc.inputParam)
			err := or.SaveClientIDSecret(context.Background(), sqlxDB, tc.inputParam)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOAuthRepository_SaveToken(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	or := NewOAuthRepository(clock.FixedClocker{})
	tests := map[string]struct {
		inputCredential *entity.MessageAPICredential
		mockSetup       func(*entity.MessageAPICredential)
		wantErr         bool
	}{
		"DB error": {
			inputCredential: &entity.MessageAPICredential{
				UserID:       1,
				AccessToken:  "OLD_ACCESS",
				RefreshToken: "OLD_REFRESH",
				ExpiresAt:    clock.FixedClocker{}.Now(),
				UpdatedAt:    clock.FixedClocker{}.Now(),
			},
			mockSetup: func(param *entity.MessageAPICredential) {
				mock.ExpectExec(`^UPDATE message_api_credentials SET access_token = \?, refresh_token = \?, expires_at = \?, updated_at = \? WHERE user_id = \?;$`).
					WithArgs(
						param.AccessToken,
						param.RefreshToken,
						param.ExpiresAt,
						param.UpdatedAt,
						param.UserID,
					).
					WillReturnError(assertAnError())
			},
			wantErr: true,
		},
		"Success": {
			inputCredential: &entity.MessageAPICredential{
				UserID:       1,
				AccessToken:  "NEW_ACCESS",
				RefreshToken: "NEW_REFRESH",
				ExpiresAt:    clock.FixedClocker{}.Now(),
				UpdatedAt:    clock.FixedClocker{}.Now(),
			},
			mockSetup: func(param *entity.MessageAPICredential) {
				mock.ExpectExec(`^UPDATE message_api_credentials SET access_token = \?, refresh_token = \?, expires_at = \?, updated_at = \? WHERE user_id = \?;$`).
					WithArgs(
						param.AccessToken,
						param.RefreshToken,
						param.ExpiresAt,
						param.UpdatedAt,
						param.UserID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup(tc.inputCredential)
			err := or.SaveToken(context.Background(), sqlxDB, tc.inputCredential)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
