package handler

import (
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

func TestDeleteMessage_ServeHTTP(t *testing.T) {
	v := validator.New()

	t.Run("ID parse error", func(t *testing.T) {
		t.Parallel()
		dm := NewDeleteMessage(&DeleteMessageServiceMock{}, v)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "abc")
		r := httptest.NewRequest(http.MethodDelete, "/messages/abc", nil)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
		w := httptest.NewRecorder()
		dm.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Equal(t, "ID must be a number", errResp.Message)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service returns ServiceError", func(t *testing.T) {
		t.Parallel()
		moq := &DeleteMessageServiceMock{
			DeleteMessageFunc: func(ctx context.Context, id entity.MessageID) error {
				return NewServiceError(
					http.StatusInternalServerError,
					"some service error",
					"detail info",
				)
			},
		}
		dm := NewDeleteMessage(moq, v)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		r := httptest.NewRequest(http.MethodDelete, "/messages/1", nil)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
		w := httptest.NewRecorder()
		dm.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Equal(t, "some service error", errResp.Message)
		assert.Equal(t, "detail info", errResp.Detail)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("service returns normal error", func(t *testing.T) {
		t.Parallel()
		moq := &DeleteMessageServiceMock{
			DeleteMessageFunc: func(ctx context.Context, id entity.MessageID) error {
				return errors.New("unexpected error")
			},
		}
		dm := NewDeleteMessage(moq, v)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		r := httptest.NewRequest(http.MethodDelete, "/messages/1", nil)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
		w := httptest.NewRecorder()
		dm.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Equal(t, "unexpected error", errResp.Message)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		moq := &DeleteMessageServiceMock{
			DeleteMessageFunc: func(ctx context.Context, id entity.MessageID) error {
				return nil
			},
		}
		dm := NewDeleteMessage(moq, v)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		r := httptest.NewRequest(http.MethodDelete, "/messages/1", nil)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
		w := httptest.NewRecorder()
		dm.ServeHTTP(w, r)
		var successResp SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &successResp)
		assert.NoError(t, err)
		assert.Equal(t, "delete message was successful", successResp.Message)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
