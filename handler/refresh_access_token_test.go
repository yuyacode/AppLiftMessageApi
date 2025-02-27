package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestRefreshAccessToken_ServeHTTP(t *testing.T) {
	v := validator.New()

	t.Run("JSON decode error", func(t *testing.T) {
		t.Parallel()
		rat := NewRefreshAccessToken(&RefreshAccessTokenServiceMock{}, v)
		r := httptest.NewRequest(http.MethodPost, "/messages/token", bytes.NewBufferString("{ invalid json }"))
		w := httptest.NewRecorder()
		rat.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Contains(t, errResp.Message, "invalid")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()
		rat := NewRefreshAccessToken(&RefreshAccessTokenServiceMock{}, v)
		body, _ := json.Marshal(map[string]string{})
		r := httptest.NewRequest(http.MethodPost, "/messages/token", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		rat.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Contains(t, errResp.Message, "required")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service returns ServiceError", func(t *testing.T) {
		t.Parallel()
		moq := &RefreshAccessTokenServiceMock{
			RefreshAccessTokenFunc: func(ctx context.Context, client_id, client_secret string) (string, string, error) {
				return "", "", NewServiceError(
					http.StatusInternalServerError,
					"invalid credentials",
					"please check your client_id or client_secret",
				)
			},
		}
		rat := NewRefreshAccessToken(moq, v)
		body, _ := json.Marshal(map[string]string{
			"client_id":     "abc",
			"client_secret": "xyz",
		})
		r := httptest.NewRequest(http.MethodPost, "/messages/token", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		rat.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Equal(t, "invalid credentials", errResp.Message)
		assert.Equal(t, "please check your client_id or client_secret", errResp.Detail)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("service returns normal error", func(t *testing.T) {
		t.Parallel()
		moq := &RefreshAccessTokenServiceMock{
			RefreshAccessTokenFunc: func(ctx context.Context, client_id, client_secret string) (string, string, error) {
				return "", "", errors.New("unexpected error")
			},
		}
		rat := NewRefreshAccessToken(moq, v)
		body, _ := json.Marshal(map[string]string{
			"client_id":     "abc",
			"client_secret": "xyz",
		})
		r := httptest.NewRequest(http.MethodPost, "/messages/token", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		rat.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Equal(t, "unexpected error", errResp.Message)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		moq := &RefreshAccessTokenServiceMock{
			RefreshAccessTokenFunc: func(ctx context.Context, client_id, client_secret string) (string, string, error) {
				return "access-token-123", "refresh-token-456", nil
			},
		}
		rat := NewRefreshAccessToken(moq, v)
		body, _ := json.Marshal(map[string]string{
			"client_id":     "abc",
			"client_secret": "xyz",
		})
		r := httptest.NewRequest(http.MethodPost, "/messages/token", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		rat.ServeHTTP(w, r)
		var rsp struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &rsp)
		assert.NoError(t, err)
		assert.Equal(t, "access-token-123", rsp.AccessToken)
		assert.Equal(t, "refresh-token-456", rsp.RefreshToken)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
