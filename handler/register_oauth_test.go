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

func TestRegisterOAuth_ServeHTTP(t *testing.T) {
	v := validator.New()

	t.Run("missing Authorization header", func(t *testing.T) {
		t.Parallel()
		ro := NewRegisterOAuth(&RegisterOAuthServiceMock{}, v)
		r := httptest.NewRequest(http.MethodPost, "/messages/register", nil)
		w := httptest.NewRecorder()
		ro.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_api_key", errResp.Message)
		assert.Equal(t, "missing Authorization header", errResp.Detail)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid Authorization header format", func(t *testing.T) {
		t.Parallel()
		ro := NewRegisterOAuth(&RegisterOAuthServiceMock{}, v)
		r := httptest.NewRequest(http.MethodPost, "/messages/register", nil)
		r.Header.Set("Authorization", "Token abc123")
		w := httptest.NewRecorder()
		ro.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_api_key", errResp.Message)
		assert.Equal(t, "invalid Authorization header format", errResp.Detail)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("empty Authorization header", func(t *testing.T) {
		t.Parallel()
		ro := NewRegisterOAuth(&RegisterOAuthServiceMock{}, v)
		r := httptest.NewRequest(http.MethodPost, "/messages/register", nil)
		r.Header.Set("Authorization", "Bearer ")
		w := httptest.NewRecorder()
		ro.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_api_key", errResp.Message)
		assert.Equal(t, "empty Authorization header", errResp.Detail)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("JSON decode error", func(t *testing.T) {
		t.Parallel()
		ro := NewRegisterOAuth(&RegisterOAuthServiceMock{}, v)
		r := httptest.NewRequest(http.MethodPost, "/messages/register", bytes.NewBufferString("{ invalid json }"))
		r.Header.Set("Authorization", "Bearer abc123")
		w := httptest.NewRecorder()
		ro.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Contains(t, errResp.Message, "invalid")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()
		ro := NewRegisterOAuth(&RegisterOAuthServiceMock{}, v)
		body, _ := json.Marshal(map[string]interface{}{
			"user_id":  123,
			"app_kind": "school",
		})
		r := httptest.NewRequest(http.MethodPost, "/messages/register", bytes.NewBuffer(body))
		r.Header.Set("Authorization", "Bearer abc123")
		w := httptest.NewRecorder()
		ro.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Contains(t, errResp.Message, "failed on the 'oneof' tag")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service returns ServiceError", func(t *testing.T) {
		t.Parallel()
		moq := &RegisterOAuthServiceMock{
			RegisterOAuthFunc: func(ctx context.Context, apiKey string) error {
				return NewServiceError(
					http.StatusInternalServerError,
					"forbidden operation",
					"cannot register OAuth",
				)
			},
		}
		ro := NewRegisterOAuth(moq, v)
		body, _ := json.Marshal(map[string]interface{}{
			"user_id":  123,
			"app_kind": "company",
		})
		r := httptest.NewRequest(http.MethodPost, "/messages/register", bytes.NewBuffer(body))
		r.Header.Set("Authorization", "Bearer abc123")
		w := httptest.NewRecorder()
		ro.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Equal(t, "forbidden operation", errResp.Message)
		assert.Equal(t, "cannot register OAuth", errResp.Detail)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("service returns normal error", func(t *testing.T) {
		t.Parallel()
		moq := &RegisterOAuthServiceMock{
			RegisterOAuthFunc: func(ctx context.Context, apiKey string) error {
				return errors.New("unexpected error")
			},
		}
		ro := NewRegisterOAuth(moq, v)
		body, _ := json.Marshal(map[string]interface{}{
			"user_id":  123,
			"app_kind": "student",
		})
		r := httptest.NewRequest(http.MethodPost, "/messages/register", bytes.NewBuffer(body))
		r.Header.Set("Authorization", "Bearer abc123")
		w := httptest.NewRecorder()
		ro.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Equal(t, "unexpected error", errResp.Message)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		moq := &RegisterOAuthServiceMock{
			RegisterOAuthFunc: func(ctx context.Context, apiKey string) error {
				return nil
			},
		}
		ro := NewRegisterOAuth(moq, v)
		body, _ := json.Marshal(map[string]interface{}{
			"user_id":  123,
			"app_kind": "company",
		})
		r := httptest.NewRequest(http.MethodPost, "/messages/register", bytes.NewBuffer(body))
		r.Header.Set("Authorization", "Bearer abc123")
		w := httptest.NewRecorder()
		ro.ServeHTTP(w, r)
		var successResp SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &successResp)
		assert.NoError(t, err)
		assert.Equal(t, "OAuth registration was successful", successResp.Message)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
