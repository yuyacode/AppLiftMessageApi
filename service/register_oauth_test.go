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

func TestRegisterOAuth_RegisterOAuth(t *testing.T) {
	type testCase struct {
		name          string
		appKind       string
		userID        int64
		apiKey        string
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
			apiKey:        "API_KEY",
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get app kind",
		},
		{
			name:    "fail to get API Key from DB",
			appKind: "company",
			userID:  1,
			apiKey:  "API_KEY",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetAPIKeyFunc = func(ctx context.Context, db store.Queryer) (string, error) {
					return "", errors.New("db error")
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get API Key",
		},
		{
			name:    "invalid API key => unauthorized",
			appKind: "company",
			userID:  1,
			apiKey:  "d9c80cbc02151d295d55f0718a18481b75a9ad604b3900b32c8d2181614c62df",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetAPIKeyFunc = func(ctx context.Context, db store.Queryer) (string, error) {
					return "137c564b6d5ff9ed412c3bd7f6e0b5d74689eac9253524e1a7d659c7ce7d59e8", nil
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusUnauthorized,
			wantErrMsg:    "API Key is invalid",
		},
		{
			name:    "client_id always exists => fail after 5 tries",
			appKind: "company",
			userID:  1,
			apiKey:  "8c967495cf41535ed0006a117f27c6a4dcb502591a6be8d600031f3c2232b77c",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetAPIKeyFunc = func(ctx context.Context, db store.Queryer) (string, error) {
					return "137c564b6d5ff9ed412c3bd7f6e0b5d74689eac9253524e1a7d659c7ce7d59e8", nil
				}
				m.SearchByClientIDFunc = func(ctx context.Context, db store.Queryer, clientID string) (bool, error) {
					return true, nil
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to generate client_id 5 times",
		},
		{
			name:    "fail to search client_id => internal server error",
			appKind: "company",
			userID:  1,
			apiKey:  "8c967495cf41535ed0006a117f27c6a4dcb502591a6be8d600031f3c2232b77c",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetAPIKeyFunc = func(ctx context.Context, db store.Queryer) (string, error) {
					return "137c564b6d5ff9ed412c3bd7f6e0b5d74689eac9253524e1a7d659c7ce7d59e8", nil
				}
				m.SearchByClientIDFunc = func(ctx context.Context, db store.Queryer, clientID string) (bool, error) {
					return false, errors.New("db error searching client_id")
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to search client_id",
		},
		{
			name:    "client_secret always exists => fail after 5 tries",
			appKind: "company",
			userID:  1,
			apiKey:  "8c967495cf41535ed0006a117f27c6a4dcb502591a6be8d600031f3c2232b77c",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetAPIKeyFunc = func(ctx context.Context, db store.Queryer) (string, error) {
					return "137c564b6d5ff9ed412c3bd7f6e0b5d74689eac9253524e1a7d659c7ce7d59e8", nil
				}
				m.SearchByClientIDFunc = func(ctx context.Context, db store.Queryer, clientID string) (bool, error) {
					return false, nil
				}
				m.SearchByClientSecretFunc = func(ctx context.Context, db store.Queryer, clientSecret string) (bool, error) {
					return true, nil
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to generate client_secret 5 times",
		},
		{
			name:    "fail to search client_secret => internal server error",
			appKind: "student",
			userID:  1,
			apiKey:  "8c967495cf41535ed0006a117f27c6a4dcb502591a6be8d600031f3c2232b77c",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetAPIKeyFunc = func(ctx context.Context, db store.Queryer) (string, error) {
					return "137c564b6d5ff9ed412c3bd7f6e0b5d74689eac9253524e1a7d659c7ce7d59e8", nil
				}
				m.SearchByClientIDFunc = func(ctx context.Context, db store.Queryer, clientID string) (bool, error) {
					return false, nil
				}
				m.SearchByClientSecretFunc = func(ctx context.Context, db store.Queryer, clientSecret string) (bool, error) {
					return false, errors.New("db error searching client_secret")
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to search client_secret",
		},
		{
			name:    "no userID in context => internal server error",
			appKind: "company",
			apiKey:  "8c967495cf41535ed0006a117f27c6a4dcb502591a6be8d600031f3c2232b77c",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetAPIKeyFunc = func(ctx context.Context, db store.Queryer) (string, error) {
					return "137c564b6d5ff9ed412c3bd7f6e0b5d74689eac9253524e1a7d659c7ce7d59e8", nil
				}
				m.SearchByClientIDFunc = func(ctx context.Context, db store.Queryer, clientID string) (bool, error) {
					return false, nil
				}
				m.SearchByClientSecretFunc = func(ctx context.Context, db store.Queryer, clientSecret string) (bool, error) {
					return false, nil
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get userID",
		},
		{
			name:    "fail SaveClientIDSecret => internal server error",
			appKind: "company",
			userID:  1,
			apiKey:  "8c967495cf41535ed0006a117f27c6a4dcb502591a6be8d600031f3c2232b77c",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetAPIKeyFunc = func(ctx context.Context, db store.Queryer) (string, error) {
					return "137c564b6d5ff9ed412c3bd7f6e0b5d74689eac9253524e1a7d659c7ce7d59e8", nil
				}
				m.SearchByClientIDFunc = func(ctx context.Context, db store.Queryer, clientID string) (bool, error) {
					return false, nil
				}
				m.SearchByClientSecretFunc = func(ctx context.Context, db store.Queryer, clientSecret string) (bool, error) {
					return false, nil
				}
			},
			prepareSetter: func(m *CredentialSetterMock) {
				m.SaveClientIDSecretFunc = func(ctx context.Context, db store.Execer, param *entity.MessageAPICredential) error {
					return errors.New("db insert error")
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to insert message api client_id and client_secret",
		},
		{
			name:    "fail searching access_token => internal server error",
			appKind: "company",
			userID:  1,
			apiKey:  "8c967495cf41535ed0006a117f27c6a4dcb502591a6be8d600031f3c2232b77c",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetAPIKeyFunc = func(ctx context.Context, db store.Queryer) (string, error) {
					return "137c564b6d5ff9ed412c3bd7f6e0b5d74689eac9253524e1a7d659c7ce7d59e8", nil
				}
				m.SearchByClientIDFunc = func(ctx context.Context, db store.Queryer, clientID string) (bool, error) {
					return false, nil
				}
				m.SearchByClientSecretFunc = func(ctx context.Context, db store.Queryer, clientSecret string) (bool, error) {
					return false, nil
				}
				m.SearchByAccessTokenFunc = func(ctx context.Context, db store.Queryer, accessToken string) (bool, error) {
					return false, errors.New("db error searching access token")
				}
			},
			prepareSetter: func(m *CredentialSetterMock) {
				m.SaveClientIDSecretFunc = func(ctx context.Context, db store.Execer, param *entity.MessageAPICredential) error {
					return nil
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to search access_token",
		},
		{
			name:    "access_token always exists => fail after 5 tries",
			appKind: "company",
			userID:  1,
			apiKey:  "8c967495cf41535ed0006a117f27c6a4dcb502591a6be8d600031f3c2232b77c",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetAPIKeyFunc = func(ctx context.Context, db store.Queryer) (string, error) {
					return "137c564b6d5ff9ed412c3bd7f6e0b5d74689eac9253524e1a7d659c7ce7d59e8", nil
				}
				m.SearchByClientIDFunc = func(ctx context.Context, db store.Queryer, clientID string) (bool, error) {
					return false, nil
				}
				m.SearchByClientSecretFunc = func(ctx context.Context, db store.Queryer, clientSecret string) (bool, error) {
					return false, nil
				}
				m.SearchByAccessTokenFunc = func(ctx context.Context, db store.Queryer, accessToken string) (bool, error) {
					return true, nil
				}
			},
			prepareSetter: func(m *CredentialSetterMock) {
				m.SaveClientIDSecretFunc = func(ctx context.Context, db store.Execer, param *entity.MessageAPICredential) error {
					return nil
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to generate access_token 5 times",
		},
		{
			name:    "refresh_token always exists => fail after 5 tries",
			appKind: "company",
			userID:  1,
			apiKey:  "8c967495cf41535ed0006a117f27c6a4dcb502591a6be8d600031f3c2232b77c",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetAPIKeyFunc = func(ctx context.Context, db store.Queryer) (string, error) {
					return "137c564b6d5ff9ed412c3bd7f6e0b5d74689eac9253524e1a7d659c7ce7d59e8", nil
				}
				m.SearchByClientIDFunc = func(ctx context.Context, db store.Queryer, clientID string) (bool, error) {
					return false, nil
				}
				m.SearchByClientSecretFunc = func(ctx context.Context, db store.Queryer, clientSecret string) (bool, error) {
					return false, nil
				}
				m.SearchByAccessTokenFunc = func(ctx context.Context, db store.Queryer, accessToken string) (bool, error) {
					return false, nil
				}
				m.SearchByRefreshTokenFunc = func(ctx context.Context, db store.Queryer, refreshToken string) (bool, error) {
					return true, nil
				}
			},
			prepareSetter: func(m *CredentialSetterMock) {
				m.SaveClientIDSecretFunc = func(ctx context.Context, db store.Execer, param *entity.MessageAPICredential) error {
					return nil
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to generate refresh_token 5 times",
		},
		{
			name:    "fail searching refresh_token => internal server error",
			appKind: "company",
			userID:  1,
			apiKey:  "8c967495cf41535ed0006a117f27c6a4dcb502591a6be8d600031f3c2232b77c",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetAPIKeyFunc = func(ctx context.Context, db store.Queryer) (string, error) {
					return "137c564b6d5ff9ed412c3bd7f6e0b5d74689eac9253524e1a7d659c7ce7d59e8", nil
				}
				m.SearchByClientIDFunc = func(ctx context.Context, db store.Queryer, clientID string) (bool, error) {
					return false, nil
				}
				m.SearchByClientSecretFunc = func(ctx context.Context, db store.Queryer, clientSecret string) (bool, error) {
					return false, nil
				}
				m.SearchByAccessTokenFunc = func(ctx context.Context, db store.Queryer, accessToken string) (bool, error) {
					return false, nil
				}
				m.SearchByRefreshTokenFunc = func(ctx context.Context, db store.Queryer, refreshToken string) (bool, error) {
					return false, errors.New("db error searching refresh_token")
				}
			},
			prepareSetter: func(m *CredentialSetterMock) {
				m.SaveClientIDSecretFunc = func(ctx context.Context, db store.Execer, param *entity.MessageAPICredential) error {
					return nil
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to search refresh_token",
		},
		{
			name:    "fail SaveToken => internal server error",
			appKind: "company",
			userID:  1,
			apiKey:  "8c967495cf41535ed0006a117f27c6a4dcb502591a6be8d600031f3c2232b77c",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetAPIKeyFunc = func(ctx context.Context, db store.Queryer) (string, error) {
					return "137c564b6d5ff9ed412c3bd7f6e0b5d74689eac9253524e1a7d659c7ce7d59e8", nil
				}
				m.SearchByClientIDFunc = func(ctx context.Context, db store.Queryer, clientID string) (bool, error) {
					return false, nil
				}
				m.SearchByClientSecretFunc = func(ctx context.Context, db store.Queryer, clientSecret string) (bool, error) {
					return false, nil
				}
				m.SearchByAccessTokenFunc = func(ctx context.Context, db store.Queryer, accessToken string) (bool, error) {
					return false, nil
				}
				m.SearchByRefreshTokenFunc = func(ctx context.Context, db store.Queryer, refreshToken string) (bool, error) {
					return false, nil
				}
			},
			prepareSetter: func(m *CredentialSetterMock) {
				m.SaveClientIDSecretFunc = func(ctx context.Context, db store.Execer, param *entity.MessageAPICredential) error {
					return nil
				}
				m.SaveTokenFunc = func(ctx context.Context, db store.Execer, param *entity.MessageAPICredential) error {
					return errors.New("db error saving token")
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to save token",
		},
		{
			name:    "success",
			appKind: "company",
			userID:  1,
			apiKey:  "8c967495cf41535ed0006a117f27c6a4dcb502591a6be8d600031f3c2232b77c",
			prepareGetter: func(m *CredentialGetterMock) {
				m.GetAPIKeyFunc = func(ctx context.Context, db store.Queryer) (string, error) {
					return "137c564b6d5ff9ed412c3bd7f6e0b5d74689eac9253524e1a7d659c7ce7d59e8", nil
				}
				m.SearchByClientIDFunc = func(ctx context.Context, db store.Queryer, clientID string) (bool, error) {
					return false, nil
				}
				m.SearchByClientSecretFunc = func(ctx context.Context, db store.Queryer, clientSecret string) (bool, error) {
					return false, nil
				}
				m.SearchByAccessTokenFunc = func(ctx context.Context, db store.Queryer, accessToken string) (bool, error) {
					return false, nil
				}
				m.SearchByRefreshTokenFunc = func(ctx context.Context, db store.Queryer, refreshToken string) (bool, error) {
					return false, nil
				}
			},
			prepareSetter: func(m *CredentialSetterMock) {
				m.SaveClientIDSecretFunc = func(ctx context.Context, db store.Execer, param *entity.MessageAPICredential) error {
					return nil
				}
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
			svc := NewRegisterOAuth(dbHandlers, getterMock, setterMock)
			err := svc.RegisterOAuth(ctx, tc.apiKey)
			if tc.wantErr {
				assert.Error(t, err, "error is expected but got nil")
				se, ok := err.(*handler.ServiceError)
				if assert.True(t, ok, "error should be *handler.ServiceError") {
					assert.Equal(t, tc.wantErrStatus, se.StatusCode)
					assert.Contains(t, se.Message, tc.wantErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
