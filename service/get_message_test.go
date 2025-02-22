package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/yuyacode/AppLiftMessageApi/entity"
	"github.com/yuyacode/AppLiftMessageApi/handler"
	"github.com/yuyacode/AppLiftMessageApi/request"
	"github.com/yuyacode/AppLiftMessageApi/store"
)

func TestGetMessage_GetAllMessages(t *testing.T) {
	type testCase struct {
		name              string
		appKind           string
		userID            int64
		prepareOwnerMock  func(*MessageOwnerGetterMock)
		prepareGetterMock func(*MessageGetterMock)
		messageThreadID   entity.MessageThreadID
		wantMessages      entity.Messages
		wantErr           bool
		wantErrStatus     int
		wantErrMsg        string
	}
	tests := []testCase{
		{
			name:          "fail if no appKind in context",
			appKind:       "",
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get app kind",
		},
		{
			name:          "fail if no userID in context",
			appKind:       "company",
			userID:        0,
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get userID",
		},
		{
			name:    "company: fail to get thread owner",
			appKind: "company",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadCompanyOwnerFunc = func(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (int64, error) {
					return 0, errors.New("owner query error")
				}
			},
			messageThreadID: 1,
			wantErr:         true,
			wantErrStatus:   http.StatusInternalServerError,
			wantErrMsg:      "failed to get threadCompanyOwner",
		},
		{
			name:    "company: user is not thread owner => forbidden",
			appKind: "company",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadCompanyOwnerFunc = func(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (int64, error) {
					return 2, nil
				}
			},
			messageThreadID: 2,
			wantErr:         true,
			wantErrStatus:   http.StatusForbidden,
			wantErrMsg:      "unauthorized: lack the necessary permissions to retrieve messages",
		},
		{
			name:    "company: fail to get messages",
			appKind: "company",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadCompanyOwnerFunc = func(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (int64, error) {
					return 1, nil
				}
			},
			prepareGetterMock: func(m *MessageGetterMock) {
				m.GetAllMessagesForCompanyUserFunc = func(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (entity.Messages, error) {
					return nil, errors.New("get messages error")
				}
			},
			messageThreadID: 1,
			wantErr:         true,
			wantErrStatus:   http.StatusInternalServerError,
			wantErrMsg:      "failed to get message",
		},
		{
			name:    "company: success",
			appKind: "company",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadCompanyOwnerFunc = func(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (int64, error) {
					return 1, nil
				}
			},
			prepareGetterMock: func(m *MessageGetterMock) {
				m.GetAllMessagesForCompanyUserFunc = func(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (entity.Messages, error) {
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
				}
			},
			messageThreadID: 1,
			wantMessages: entity.Messages{
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
			},
			wantErr: false,
		},
		{
			name:    "student: fail to get thread owner",
			appKind: "student",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadStudentOwnerFunc = func(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (int64, error) {
					return 0, errors.New("owner query error")
				}
			},
			messageThreadID: 1,
			wantErr:         true,
			wantErrStatus:   http.StatusInternalServerError,
			wantErrMsg:      "failed to get threadStudentOwner",
		},
		{
			name:    "student: user is not thread owner => forbidden",
			appKind: "student",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadStudentOwnerFunc = func(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (int64, error) {
					return 2, nil
				}
			},
			messageThreadID: 2,
			wantErr:         true,
			wantErrStatus:   http.StatusForbidden,
			wantErrMsg:      "unauthorized: lack the necessary permissions to retrieve messages",
		},
		{
			name:    "student: fail to get messages",
			appKind: "student",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadStudentOwnerFunc = func(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (int64, error) {
					return 1, nil
				}
			},
			prepareGetterMock: func(m *MessageGetterMock) {
				m.GetAllMessagesForStudentUserFunc = func(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (entity.Messages, error) {
					return nil, errors.New("get messages error")
				}
			},
			messageThreadID: 1,
			wantErr:         true,
			wantErrStatus:   http.StatusInternalServerError,
			wantErrMsg:      "failed to get message",
		},
		{
			name:    "student: success",
			appKind: "student",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadStudentOwnerFunc = func(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (int64, error) {
					return 1, nil
				}
			},
			prepareGetterMock: func(m *MessageGetterMock) {
				m.GetAllMessagesForStudentUserFunc = func(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (entity.Messages, error) {
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
				}
			},
			messageThreadID: 1,
			wantMessages: entity.Messages{
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
			},
			wantErr: false,
		},
	}
	dbHandlers := map[string]*sqlx.DB{
		"common": nil,
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			if tc.appKind != "" {
				ctx = request.SetAppKind(ctx, tc.appKind)
			}
			if tc.userID != 0 {
				ctx = request.SetUserID(ctx, tc.userID)
			}
			ownerMock := &MessageOwnerGetterMock{}
			getterMock := &MessageGetterMock{}
			if tc.prepareOwnerMock != nil {
				tc.prepareOwnerMock(ownerMock)
			}
			if tc.prepareGetterMock != nil {
				tc.prepareGetterMock(getterMock)
			}
			svc := NewGetMessage(dbHandlers, getterMock, ownerMock)
			messages, err := svc.GetAllMessages(ctx, tc.messageThreadID)
			if tc.wantErr {
				assert.Error(t, err, "error is expected but got nil")
				se, ok := err.(*handler.ServiceError)
				if assert.True(t, ok, "error should be *handler.ServiceError") {
					assert.Equal(t, tc.wantErrStatus, se.StatusCode)
					assert.Contains(t, se.Message, tc.wantErrMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantMessages, messages)
			}
		})
	}
}
