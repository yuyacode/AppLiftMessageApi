package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/yuyacode/AppLiftMessageApi/entity"
	"github.com/yuyacode/AppLiftMessageApi/handler"
	"github.com/yuyacode/AppLiftMessageApi/request"
	"github.com/yuyacode/AppLiftMessageApi/store"
)

func TestRefreshAccessToken_RefreshAccessToken(t *testing.T) {
	type testCase struct {
		name          string
		appKind       string
		userID        int64
		clientID      string
		clientSecret  string
		prepareGetter func(*CredentialGetterMock)
		prepareSetter func(*CredentialSetterMock)
		wantErr       bool
		wantErrStatus int
		wantErrMsg    string
	}
	tests := []testCase{
		{
			name:          "fail if no appKind in context",
			userID:        1,
			clientID:      "client123",
			clientSecret:  "secret123",
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get app kind",
		},
		{
			name:          "fail if no userID in context",
			appKind:       "company",
			clientID:      "client123",
			clientSecret:  "secret123",
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get user_id",
		},
		{
			name:         "fail to get clientID from DB",
			appKind:      "company",
			userID:       1,
			clientID:     "client123",
			clientSecret: "secret123",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetClientIDFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "", errors.New("db error")
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get client_id",
		},
		{
			name:         "client_id mismatch => unauthorized",
			appKind:      "company",
			userID:       1,
			clientID:     "bad-client",
			clientSecret: "secret123",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetClientIDFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "valid-client", nil
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusUnauthorized,
			wantErrMsg:    "client_id is invalid",
		},
		{
			name:         "fail to get clientSecret from DB",
			appKind:      "company",
			userID:       1,
			clientID:     "client123",
			clientSecret: "secret123",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetClientIDFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "client123", nil
				}
				m.GetClientSecretFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "", errors.New("db error")
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get client_secret",
		},
		{
			name:         "client_secret mismatch => unauthorized",
			appKind:      "company",
			userID:       1,
			clientID:     "client123",
			clientSecret: "wrong-secret",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetClientIDFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "client123", nil
				}
				m.GetClientSecretFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "valid-secret", nil
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusUnauthorized,
			wantErrMsg:    "client_secret is invalid",
		},
		{
			name:         "error searching for access token => internal server error",
			appKind:      "company",
			userID:       1,
			clientID:     "client123",
			clientSecret: "secret123",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetClientIDFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "client123", nil
				}
				m.GetClientSecretFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "secret123", nil
				}
				m.SearchByAccessTokenFunc = func(ctx context.Context, db store.Queryer, accessToken string) (bool, error) {
					return false, errors.New("db error searching access token")
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to search access_token",
		},
		{
			name:         "access token always exist => fail after 5 tries",
			appKind:      "company",
			userID:       1,
			clientID:     "client123",
			clientSecret: "secret123",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetClientIDFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "client123", nil
				}
				m.GetClientSecretFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "secret123", nil
				}
				m.SearchByAccessTokenFunc = func(ctx context.Context, db store.Queryer, accessToken string) (bool, error) {
					return true, nil
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to generate access_token 5 times",
		},
		{
			name:         "error searching for refresh token => internal server error",
			appKind:      "student",
			userID:       1,
			clientID:     "client-student",
			clientSecret: "secret-student",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetClientIDFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "client-student", nil
				}
				m.GetClientSecretFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "secret-student", nil
				}
				m.SearchByAccessTokenFunc = func(ctx context.Context, db store.Queryer, accessToken string) (bool, error) {
					return false, nil
				}
				m.SearchByRefreshTokenFunc = func(ctx context.Context, db store.Queryer, refreshToken string) (bool, error) {
					return false, errors.New("db error searching refresh token")
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to search refresh_token",
		},
		{
			name:         "refresh token always exist => fail after 5 tries",
			appKind:      "student",
			userID:       1,
			clientID:     "client-student",
			clientSecret: "secret-student",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetClientIDFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "client-student", nil
				}
				m.GetClientSecretFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "secret-student", nil
				}
				m.SearchByAccessTokenFunc = func(ctx context.Context, db store.Queryer, accessToken string) (bool, error) {
					return false, nil
				}
				m.SearchByRefreshTokenFunc = func(ctx context.Context, db store.Queryer, refreshToken string) (bool, error) {
					return true, nil
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to generate refresh_token 5 times",
		},
		{
			name:         "fail to save token => internal server error",
			appKind:      "company",
			userID:       1,
			clientID:     "client123",
			clientSecret: "secret123",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetClientIDFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "client123", nil
				}
				m.GetClientSecretFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "secret123", nil
				}
				m.SearchByAccessTokenFunc = func(ctx context.Context, db store.Queryer, accessToken string) (bool, error) {
					return false, nil
				}
				m.SearchByRefreshTokenFunc = func(ctx context.Context, db store.Queryer, refreshToken string) (bool, error) {
					return false, nil
				}
			},
			prepareSetter: func(m *CredentialSetterMock) {
				m.SaveTokenFunc = func(ctx context.Context, db store.Execer, param *entity.MessageAPICredential) error {
					return errors.New("db save error")
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to save token",
		},
		{
			name:         "success",
			appKind:      "company",
			userID:       1,
			clientID:     "client123",
			clientSecret: "secret123",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetClientIDFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "client123", nil
				}
				m.GetClientSecretFunc = func(ctx context.Context, db store.Queryer, userID int64) (string, error) {
					return "secret123", nil
				}
				m.SearchByAccessTokenFunc = func(ctx context.Context, db store.Queryer, accessToken string) (bool, error) {
					return false, nil
				}
				m.SearchByRefreshTokenFunc = func(ctx context.Context, db store.Queryer, refreshToken string) (bool, error) {
					return false, nil
				}
			},
			prepareSetter: func(m *CredentialSetterMock) {
				m.SaveTokenFunc = func(ctx context.Context, db store.Execer, param *entity.MessageAPICredential) error {
					return nil
				}
			},
			wantErr: false,
		},
	}
	dbHandlers := map[string]*sqlx.DB{
		"company": nil,
		"student": nil,
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			if tc.appKind != "" {
				ctx = request.SetAppKind(ctx, tc.appKind)
			}
			if tc.userID != 0 {
				ctx = request.SetUserID(ctx, tc.userID)
			}
			getterMock := &CredentialGetterMock{}
			setterMock := &CredentialSetterMock{}
			if tc.prepareGetter != nil {
				tc.prepareGetter(getterMock)
			}
			if tc.prepareSetter != nil {
				tc.prepareSetter(setterMock)
			}
			svc := NewRefreshAccessToken(dbHandlers, getterMock, setterMock)
			accessToken, refreshToken, err := svc.RefreshAccessToken(ctx, tc.clientID, tc.clientSecret)
			if tc.wantErr {
				assert.Error(t, err, "error is expected but got nil")
				se, ok := err.(*handler.ServiceError)
				if assert.True(t, ok, "error should be *handler.ServiceError") {
					assert.Equal(t, tc.wantErrStatus, se.StatusCode)
					assert.Contains(t, se.Message, tc.wantErrMsg)
				}
				assert.Empty(t, accessToken)
				assert.Empty(t, refreshToken)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, accessToken)
				assert.NotEmpty(t, refreshToken)
			}
		})
	}
}
