package store

import (
	"context"
	"testing"

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
