package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/yuyacode/AppLiftMessageApi/entity"
)

func TestAddMessage_ServeHTTP(t *testing.T) {
	v := validator.New()

	t.Run("JSON decode error", func(t *testing.T) {
		t.Parallel()
		am := NewAddMessage(&AddMessageServiceMock{}, v)
		r := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewBufferString("{ invalid json }"))
		w := httptest.NewRecorder()
		am.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Contains(t, errResp.Message, "invalid")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()
		am := NewAddMessage(&AddMessageServiceMock{}, v)
		requestBody, _ := json.Marshal(map[string]interface{}{
			// "message_thread_id" を省略
			"is_from_company": 1,
			"is_from_student": 0,
			"content":         "Hello",
			"is_sent":         1,
			"sent_at":         time.Now().Format(time.RFC3339),
		})
		r := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewBuffer(requestBody))
		w := httptest.NewRecorder()
		am.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Contains(t, errResp.Message, "required")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service returns ServiceError", func(t *testing.T) {
		t.Parallel()
		moq := &AddMessageServiceMock{
			AddMessageFunc: func(ctx context.Context, messageThreadID entity.MessageThreadID, isFromCompany int8, isFromStudent int8, content string, isSent int8, sentAt time.Time) (*entity.Message, error) {
				return nil, NewServiceError(
					http.StatusInternalServerError,
					"some service error",
					"something detail",
				)
			},
		}
		am := NewAddMessage(moq, v)
		requestBody, _ := json.Marshal(map[string]interface{}{
			"message_thread_id": 1,
			"is_from_company":   1,
			"is_from_student":   0,
			"content":           "Hello",
			"is_sent":           1,
			"sent_at":           time.Now().Format(time.RFC3339),
		})
		r := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewBuffer(requestBody))
		w := httptest.NewRecorder()
		am.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Equal(t, "some service error", errResp.Message)
		assert.Equal(t, "something detail", errResp.Detail)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("service returns normal error", func(t *testing.T) {
		t.Parallel()
		moq := &AddMessageServiceMock{
			AddMessageFunc: func(ctx context.Context, messageThreadID entity.MessageThreadID, isFromCompany int8, isFromStudent int8, content string, isSent int8, sentAt time.Time) (*entity.Message, error) {
				return nil, errors.New("unexpected error")
			},
		}
		am := NewAddMessage(moq, v)
		requestBody, _ := json.Marshal(map[string]interface{}{
			"message_thread_id": 1,
			"is_from_company":   1,
			"is_from_student":   0,
			"content":           "Hello",
			"is_sent":           1,
			"sent_at":           time.Now().Format(time.RFC3339),
		})
		r := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewBuffer(requestBody))
		w := httptest.NewRecorder()
		am.ServeHTTP(w, r)
		var errResp ErrResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Equal(t, "unexpected error", errResp.Message)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		moq := &AddMessageServiceMock{
			AddMessageFunc: func(ctx context.Context, messageThreadID entity.MessageThreadID, isFromCompany int8, isFromStudent int8, content string, isSent int8, sentAt time.Time) (*entity.Message, error) {
				return &entity.Message{
					ID:              entity.MessageID(1),
					MessageThreadID: messageThreadID,
					IsFromCompany:   isFromCompany,
					IsFromStudent:   isFromStudent,
					Content:         content,
					IsSent:          isSent,
					SentAt:          sentAt,
				}, nil
			},
		}
		am := NewAddMessage(moq, v)
		requestBody, _ := json.Marshal(map[string]interface{}{
			"message_thread_id": 1,
			"is_from_company":   1,
			"is_from_student":   0,
			"content":           "Hello",
			"is_sent":           1,
			"sent_at":           time.Now().Format(time.RFC3339),
		})
		r := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewBuffer(requestBody))
		w := httptest.NewRecorder()
		am.ServeHTTP(w, r)
		var rsp struct {
			ID entity.MessageID `json:"id"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &rsp)
		assert.NoError(t, err)
		assert.Equal(t, entity.MessageID(1), rsp.ID)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
