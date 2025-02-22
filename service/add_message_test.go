package service

import (
	"context"
	"database/sql"
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

func TestAddMessage_AddMessage(t *testing.T) {
	type testCase struct {
		name             string
		appKind          string
		userID           int64
		prepareOwnerMock func(*MessageOwnerGetterMock)
		prepareAdderMock func(*MessageAdderMock)
		messageThreadID  entity.MessageThreadID
		isFromCompany    int8
		isFromStudent    int8
		content          string
		isSent           int8
		sentAt           time.Time
		wantMsgID        entity.MessageID
		wantErr          bool
		wantErrStatus    int
		wantErrMsg       string
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
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get threadCompanyOwner",
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
			wantErr:       true,
			wantErrStatus: http.StatusForbidden,
			wantErrMsg:    "unauthorized: lack the necessary permissions to add messages",
		},
		{
			name:    "company: fail to add message",
			appKind: "company",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadCompanyOwnerFunc = func(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (int64, error) {
					return 1, nil
				}
			},
			prepareAdderMock: func(m *MessageAdderMock) {
				m.AddMessageFunc = func(ctx context.Context, db store.Execer, param *entity.Message) error {
					return errors.New("add message error")
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to add message",
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
			prepareAdderMock: func(m *MessageAdderMock) {
				m.AddMessageFunc = func(ctx context.Context, db store.Execer, param *entity.Message) error {
					param.CreatedAt = &sql.NullTime{
						Time:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					}
					param.ID = entity.MessageID(1)
					return nil
				}
			},
			messageThreadID: entity.MessageThreadID(1),
			isFromCompany:   1,
			isFromStudent:   0,
			content:         "message from company user",
			isSent:          1,
			sentAt:          time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantMsgID:       entity.MessageID(1),
			wantErr:         false,
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
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get threadStudentOwner",
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
			wantErr:       true,
			wantErrStatus: http.StatusForbidden,
			wantErrMsg:    "unauthorized: lack the necessary permissions to add messages",
		},
		{
			name:    "student: fail to add message",
			appKind: "student",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadStudentOwnerFunc = func(ctx context.Context, db store.Queryer, messageThreadID entity.MessageThreadID) (int64, error) {
					return 1, nil
				}
			},
			prepareAdderMock: func(m *MessageAdderMock) {
				m.AddMessageFunc = func(ctx context.Context, db store.Execer, param *entity.Message) error {
					return errors.New("add message error")
				}
			},
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to add message",
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
			prepareAdderMock: func(m *MessageAdderMock) {
				m.AddMessageFunc = func(ctx context.Context, db store.Execer, param *entity.Message) error {
					param.CreatedAt = &sql.NullTime{
						Time:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					}
					param.ID = entity.MessageID(1)
					return nil
				}
			},
			messageThreadID: entity.MessageThreadID(1),
			isFromCompany:   0,
			isFromStudent:   1,
			content:         "message from student user",
			isSent:          1,
			sentAt:          time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantMsgID:       entity.MessageID(1),
			wantErr:         false,
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
			adderMock := &MessageAdderMock{}
			if tc.prepareOwnerMock != nil {
				tc.prepareOwnerMock(ownerMock)
			}
			if tc.prepareAdderMock != nil {
				tc.prepareAdderMock(adderMock)
			}
			svc := NewAddMessage(dbHandlers, adderMock, ownerMock)
			msg, err := svc.AddMessage(ctx, tc.messageThreadID, tc.isFromCompany, tc.isFromStudent, tc.content, tc.isSent, tc.sentAt)
			if tc.wantErr {
				assert.Error(t, err, "error is expected but got nil")
				se, ok := err.(*handler.ServiceError)
				if assert.True(t, ok, "error should be *handler.ServiceError") {
					assert.Equal(t, tc.wantErrStatus, se.StatusCode)
					assert.Contains(t, se.Message, tc.wantErrMsg)
				}
				assert.Nil(t, msg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, msg)
				assert.Equal(t, tc.wantMsgID, msg.ID)
				assert.Equal(t, tc.messageThreadID, msg.MessageThreadID)
				assert.Equal(t, tc.isFromCompany, msg.IsFromCompany)
				assert.Equal(t, tc.isFromStudent, msg.IsFromStudent)
				assert.Equal(t, tc.content, msg.Content)
				assert.Equal(t, tc.isSent, msg.IsSent)
				assert.Equal(t, tc.sentAt, msg.SentAt)
				assert.Equal(t, &sql.NullTime{
					Time:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					Valid: true,
				}, msg.CreatedAt)
			}
		})
	}
}
