package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/yuyacode/AppLiftMessageApi/entity"
)

func TestGetMessage_ServeHTTP(t *testing.T) {
	t.Run("missing thread_id query param", func(t *testing.T) {
		t.Parallel()
		gm := &GetMessage{}
		r := httptest.NewRequest(http.MethodGet, "/messages", nil)
		w := httptest.NewRecorder()
		gm.ServeHTTP(w, r)
		var errResp ErrResponse
		json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.Equal(t, errResp.Message, "missing required query parameter: thread_id")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid thread_id format", func(t *testing.T) {
		t.Parallel()
		gm := &GetMessage{}
		r := httptest.NewRequest(http.MethodGet, "/messages?thread_id=abc", nil)
		w := httptest.NewRecorder()
		gm.ServeHTTP(w, r)
		var errResp ErrResponse
		json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.Equal(t, errResp.Message, "invalid format for query parameter: thread_id. Must be a valid integer")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("service returns ServiceError", func(t *testing.T) {
		t.Parallel()
		moq := &GetMessageServiceMock{
			GetAllMessagesFunc: func(ctx context.Context, messageThreadID entity.MessageThreadID) (entity.Messages, error) {
				return nil, NewServiceError(
					http.StatusInternalServerError,
					"some service error",
					"something detail",
				)
			},
		}
		gm := &GetMessage{Service: moq}
		r := httptest.NewRequest(http.MethodGet, "/messages?thread_id=1", nil)
		w := httptest.NewRecorder()
		gm.ServeHTTP(w, r)
		var errResp ErrResponse
		json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.Equal(t, "some service error", errResp.Message)
		assert.Equal(t, "something detail", errResp.Detail)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("service returns normal error", func(t *testing.T) {
		t.Parallel()
		moq := &GetMessageServiceMock{
			GetAllMessagesFunc: func(ctx context.Context, messageThreadID entity.MessageThreadID) (entity.Messages, error) {
				return nil, errors.New("unexpected error")
			},
		}
		gm := &GetMessage{Service: moq}
		r := httptest.NewRequest(http.MethodGet, "/messages?thread_id=1", nil)
		w := httptest.NewRecorder()
		gm.ServeHTTP(w, r)
		var errResp ErrResponse
		json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.Equal(t, "unexpected error", errResp.Message)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		moq := &GetMessageServiceMock{
			GetAllMessagesFunc: func(ctx context.Context, messageThreadID entity.MessageThreadID) (entity.Messages, error) {
				return entity.Messages{
					&entity.Message{
						ID:            entity.MessageID(1),
						IsFromCompany: 1,
						IsFromStudent: 0,
						Content:       "normal message from company user",
						IsSent:        1,
						SentAt:        time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					&entity.Message{
						ID:            entity.MessageID(2),
						IsFromCompany: 0,
						IsFromStudent: 1,
						Content:       "reservation message from student user",
						IsSent:        0,
						SentAt:        time.Date(2025, 1, 2, 9, 0, 0, 0, time.UTC),
					},
				}, nil
			},
		}
		gm := &GetMessage{Service: moq}
		r := httptest.NewRequest(http.MethodGet, "/messages?thread_id=1", nil)
		w := httptest.NewRecorder()
		gm.ServeHTTP(w, r)

		var messages []message
		json.Unmarshal(w.Body.Bytes(), &messages)
		assert.Len(t, messages, 2)

		assert.Equal(t, entity.MessageID(1), messages[0].ID)
		assert.Equal(t, int8(1), messages[0].IsFromCompany)
		assert.Equal(t, int8(0), messages[0].IsFromStudent)
		assert.Equal(t, "normal message from company user", messages[0].Content)
		assert.Equal(t, int8(1), messages[0].IsSent)
		assert.Equal(t, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), messages[0].SentAt)

		assert.Equal(t, entity.MessageID(2), messages[1].ID)
		assert.Equal(t, int8(0), messages[1].IsFromCompany)
		assert.Equal(t, int8(1), messages[1].IsFromStudent)
		assert.Equal(t, "reservation message from student user", messages[1].Content)
		assert.Equal(t, int8(0), messages[1].IsSent)
		assert.Equal(t, time.Date(2025, 1, 2, 9, 0, 0, 0, time.UTC), messages[1].SentAt)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
