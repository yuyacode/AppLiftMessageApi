package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yuyacode/AppLiftMessageApi/clock"
	"github.com/yuyacode/AppLiftMessageApi/entity"
)

type mockError struct {
	msg string
}

func (m *mockError) Error() string {
	return m.msg
}

func assertAnError() error {
	return &mockError{"some db error"}
}

func newMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Close()
	})
	sqlxDB := sqlx.NewDb(db, "mysql")
	return sqlxDB, mock
}

func TestMessageRepository_GetThreadCompanyOwner(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	mr := NewMessageRepository(clock.FixedClocker{})
	tests := map[string]struct {
		messageThreadID entity.MessageThreadID
		mockSetup       func()
		wantErr         bool
		wantResult      int64
	}{
		"DB error": {
			messageThreadID: 1,
			mockSetup: func() {
				mock.ExpectQuery("^SELECT company_user_id FROM message_threads WHERE id = \\? AND deleted_at IS NULL;$").
					WithArgs(int64(1)).
					WillReturnError(assertAnError())
			},
			wantErr:    true,
			wantResult: 0,
		},
		"Success": {
			messageThreadID: 1,
			mockSetup: func() {
				mock.ExpectQuery("^SELECT company_user_id FROM message_threads WHERE id = \\? AND deleted_at IS NULL;$").
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"company_user_id"}).AddRow(int64(9999)))
			},
			wantErr:    false,
			wantResult: 9999,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			got, err := mr.GetThreadCompanyOwner(context.Background(), sqlxDB, tc.messageThreadID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantResult, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMessageRepository_GetThreadStudentOwner(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	mr := NewMessageRepository(clock.FixedClocker{})
	tests := map[string]struct {
		messageThreadID entity.MessageThreadID
		mockSetup       func()
		wantErr         bool
		wantResult      int64
	}{
		"DB error": {
			messageThreadID: 1,
			mockSetup: func() {
				mock.ExpectQuery("^SELECT student_user_id FROM message_threads WHERE id = \\? AND deleted_at IS NULL;$").
					WithArgs(int64(1)).
					WillReturnError(assertAnError())
			},
			wantErr:    true,
			wantResult: 0,
		},
		"Success": {
			messageThreadID: 1,
			mockSetup: func() {
				mock.ExpectQuery("^SELECT student_user_id FROM message_threads WHERE id = \\? AND deleted_at IS NULL;$").
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"student_user_id"}).AddRow(int64(9999)))
			},
			wantErr:    false,
			wantResult: 9999,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			got, err := mr.GetThreadStudentOwner(context.Background(), sqlxDB, tc.messageThreadID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantResult, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
