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

func TestEditMessage_EditMessage(t *testing.T) {
	type testCase struct {
		name              string
		appKind           string
		userID            int64
		prepareOwnerMock  func(*MessageOwnerGetterMock)
		prepareEditorMock func(*MessageEditorMock)
		messageID         entity.MessageID
		content           string
		wantErr           bool
		wantErrStatus     int
		wantErrMsg        string
	}
	tests := []testCase{
		{
			name:          "fail if no appKind in context",
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get app kind",
		},
		{
			name:          "fail if no userID in context",
			appKind:       "company",
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get userID",
		},
		{
			name:    "company: fail to get thread owner by message ID",
			appKind: "company",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadCompanyOwnerByMessageIDFunc = func(ctx context.Context, db store.Queryer, messageID entity.MessageID) (int64, error) {
					return 0, errors.New("owner query error")
				}
			},
			messageID:     1,
			content:       "edited content",
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get threadCompanyOwner",
		},
		{
			name:    "company: user mismatch => forbidden",
			appKind: "company",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadCompanyOwnerByMessageIDFunc = func(ctx context.Context, db store.Queryer, messageID entity.MessageID) (int64, error) {
					return 2, nil
				}
			},
			messageID:     2,
			content:       "edited content",
			wantErr:       true,
			wantErrStatus: http.StatusForbidden,
			wantErrMsg:    "unauthorized: lack the necessary permissions to edit message",
		},
		{
			name:    "company: editor fails => internal server error",
			appKind: "company",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadCompanyOwnerByMessageIDFunc = func(ctx context.Context, db store.Queryer, messageID entity.MessageID) (int64, error) {
					return 1, nil
				}
			},
			prepareEditorMock: func(m *MessageEditorMock) {
				m.EditMessageFunc = func(ctx context.Context, db store.Execer, param *entity.Message) error {
					return errors.New("update error")
				}
			},
			messageID:     1,
			content:       "edited content",
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to edit message",
		},
		{
			name:    "company: success",
			appKind: "company",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadCompanyOwnerByMessageIDFunc = func(ctx context.Context, db store.Queryer, messageID entity.MessageID) (int64, error) {
					return 1, nil
				}
			},
			prepareEditorMock: func(m *MessageEditorMock) {
				m.EditMessageFunc = func(ctx context.Context, db store.Execer, param *entity.Message) error {
					return nil
				}
			},
			messageID: 1,
			content:   "edited content",
			wantErr:   false,
		},
		{
			name:    "student: fail to get thread owner by message ID",
			appKind: "student",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadStudentOwnerByMessageIDFunc = func(ctx context.Context, db store.Queryer, messageID entity.MessageID) (int64, error) {
					return 0, errors.New("owner query error")
				}
			},
			messageID:     1,
			content:       "edited content",
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to get threadStudentOwner",
		},
		{
			name:    "student: user mismatch => forbidden",
			appKind: "student",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadStudentOwnerByMessageIDFunc = func(ctx context.Context, db store.Queryer, messageID entity.MessageID) (int64, error) {
					return 2, nil
				}
			},
			messageID:     2,
			content:       "edited content",
			wantErr:       true,
			wantErrStatus: http.StatusForbidden,
			wantErrMsg:    "unauthorized: lack the necessary permissions to edit message",
		},
		{
			name:    "student: editor fails => internal server error",
			appKind: "student",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadStudentOwnerByMessageIDFunc = func(ctx context.Context, db store.Queryer, messageID entity.MessageID) (int64, error) {
					return 1, nil
				}
			},
			prepareEditorMock: func(m *MessageEditorMock) {
				m.EditMessageFunc = func(ctx context.Context, db store.Execer, param *entity.Message) error {
					return errors.New("update error")
				}
			},
			messageID:     1,
			content:       "edited content",
			wantErr:       true,
			wantErrStatus: http.StatusInternalServerError,
			wantErrMsg:    "failed to edit message",
		},
		{
			name:    "student: success",
			appKind: "student",
			userID:  1,
			prepareOwnerMock: func(m *MessageOwnerGetterMock) {
				m.GetThreadStudentOwnerByMessageIDFunc = func(ctx context.Context, db store.Queryer, messageID entity.MessageID) (int64, error) {
					return 1, nil
				}
			},
			prepareEditorMock: func(m *MessageEditorMock) {
				m.EditMessageFunc = func(ctx context.Context, db store.Execer, param *entity.Message) error {
					return nil
				}
			},
			messageID: 1,
			content:   "edited content",
			wantErr:   false,
		},
	}
	dbHandlers := map[string]*sqlx.DB{
		"common": nil,
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
			ownerMock := &MessageOwnerGetterMock{}
			editorMock := &MessageEditorMock{}
			if tc.prepareOwnerMock != nil {
				tc.prepareOwnerMock(ownerMock)
			}
			if tc.prepareEditorMock != nil {
				tc.prepareEditorMock(editorMock)
			}
			svc := NewEditMessage(dbHandlers, editorMock, ownerMock)
			err := svc.EditMessage(ctx, tc.messageID, tc.content)
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
