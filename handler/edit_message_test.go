package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/yuyacode/AppLiftMessageApi/entity"
)

func TestEditMessage_ServeHTTP(t *testing.T) {
	v := validator.New()

	t.Run("ID parse error", func(t *testing.T) {
		t.Parallel()
		em := NewEditMessage(&EditMessageServiceMock{}, v)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "abc")
		r := httptest.NewRequest(http.MethodPatch, "/messages/abc", nil)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
		w := httptest.NewRecorder()
		em.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Equal(t, "ID must be a number", errResp.Message)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("JSON decode error", func(t *testing.T) {
		t.Parallel()
		em := NewEditMessage(&EditMessageServiceMock{}, v)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		r := httptest.NewRequest(http.MethodPatch, "/messages/1", bytes.NewBufferString("{ invalid json }"))
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
		w := httptest.NewRecorder()
		em.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Contains(t, errResp.Message, "invalid")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()
		em := NewEditMessage(&EditMessageServiceMock{}, v)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		body, _ := json.Marshal(map[string]interface{}{})
		r := httptest.NewRequest(http.MethodPatch, "/messages/1", bytes.NewBuffer(body))
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
		w := httptest.NewRecorder()
		em.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Contains(t, errResp.Message, "required")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service returns ServiceError", func(t *testing.T) {
		t.Parallel()
		moq := &EditMessageServiceMock{
			EditMessageFunc: func(ctx context.Context, id entity.MessageID, content string) error {
				return NewServiceError(
					http.StatusInternalServerError,
					"some service error",
					"detail info",
				)
			},
		}
		em := NewEditMessage(moq, v)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		body, _ := json.Marshal(map[string]string{"content": "updated content"})
		r := httptest.NewRequest(http.MethodPatch, "/messages/1", bytes.NewBuffer(body))
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
		w := httptest.NewRecorder()
		em.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Equal(t, "some service error", errResp.Message)
		assert.Equal(t, "detail info", errResp.Detail)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("service returns normal error", func(t *testing.T) {
		t.Parallel()
		moq := &EditMessageServiceMock{
			EditMessageFunc: func(ctx context.Context, id entity.MessageID, content string) error {
				return errors.New("unexpected error")
			},
		}
		em := NewEditMessage(moq, v)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		body, _ := json.Marshal(map[string]string{"content": "updated content"})
		r := httptest.NewRequest(http.MethodPatch, "/messages/1", bytes.NewBuffer(body))
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
		w := httptest.NewRecorder()
		em.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Equal(t, "unexpected error", errResp.Message)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		moq := &EditMessageServiceMock{
			EditMessageFunc: func(ctx context.Context, id entity.MessageID, content string) error {
				return nil
			},
		}
		em := NewEditMessage(moq, v)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		body, _ := json.Marshal(map[string]string{"content": "updated content"})
		r := httptest.NewRequest(http.MethodPatch, "/messages/1", bytes.NewBuffer(body))
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
		w := httptest.NewRecorder()
		em.ServeHTTP(w, r)
		var successResp SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &successResp)
		assert.NoError(t, err)
		assert.Equal(t, "edit message was successful", successResp.Message)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
