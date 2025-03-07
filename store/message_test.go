package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yuyacode/AppLiftMessageApi/clock"
	"github.com/yuyacode/AppLiftMessageApi/entity"
)

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

func TestMessageRepository_GetThreadCompanyOwnerByMessageID(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	mr := NewMessageRepository(clock.FixedClocker{})
	tests := map[string]struct {
		messageID  entity.MessageID
		mockSetup  func()
		wantErr    bool
		wantResult int64
	}{
		"DB error": {
			messageID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT company_user_id\s+FROM message_threads\s+INNER JOIN messages\s+ON message_threads.id = messages.message_thread_id\s+WHERE messages.id = \? AND is_from_company = 1;$`).
					WithArgs(int64(1)).
					WillReturnError(assertAnError())
			},
			wantErr:    true,
			wantResult: 0,
		},
		"Success": {
			messageID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT company_user_id\s+FROM message_threads\s+INNER JOIN messages\s+ON message_threads.id = messages.message_thread_id\s+WHERE messages.id = \? AND is_from_company = 1;$`).
					WithArgs(int64(1)).
					WillReturnRows(
						sqlmock.NewRows([]string{"company_user_id"}).AddRow(int64(9999)),
					)
			},
			wantErr:    false,
			wantResult: 9999,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			got, err := mr.GetThreadCompanyOwnerByMessageID(context.Background(), sqlxDB, tc.messageID)
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

func TestMessageRepository_GetThreadStudentOwnerByMessageID(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	mr := NewMessageRepository(clock.FixedClocker{})
	tests := map[string]struct {
		messageID  entity.MessageID
		mockSetup  func()
		wantErr    bool
		wantResult int64
	}{
		"DB error": {
			messageID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT student_user_id\s+FROM message_threads\s+INNER JOIN messages\s+ON message_threads.id = messages.message_thread_id\s+WHERE messages.id = \? AND is_from_student = 1;$`).
					WithArgs(int64(1)).
					WillReturnError(assertAnError())
			},
			wantErr:    true,
			wantResult: 0,
		},
		"Success": {
			messageID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT student_user_id\s+FROM message_threads\s+INNER JOIN messages\s+ON message_threads.id = messages.message_thread_id\s+WHERE messages.id = \? AND is_from_student = 1;$`).
					WithArgs(int64(1)).
					WillReturnRows(
						sqlmock.NewRows([]string{"company_user_id"}).AddRow(int64(9999)),
					)
			},
			wantErr:    false,
			wantResult: 9999,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			got, err := mr.GetThreadStudentOwnerByMessageID(context.Background(), sqlxDB, tc.messageID)
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

func TestMessageRepository_GetAllMessagesForCompanyUser(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	mr := NewMessageRepository(clock.FixedClocker{})
	tests := map[string]struct {
		messageThreadID entity.MessageThreadID
		mockSetup       func()
		wantErr         bool
		wantMessages    entity.Messages
	}{
		"DB error": {
			messageThreadID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT id, is_from_company, is_from_student, content, is_sent, sent_at\s+FROM messages\s+WHERE message_thread_id = \?\s+AND deleted_at IS NULL\s+AND\s+\(\s*\(is_from_company = 1 AND is_sent = 0\)\s+OR \(is_from_company = 1 AND is_sent = 1\)\s+OR \(is_from_student = 1 AND is_sent = 1\)\s*\)\s+ORDER BY sent_at ASC, id ASC;$`).
					WithArgs(int64(1)).
					WillReturnError(assertAnError())
			},
			wantErr:      true,
			wantMessages: nil,
		},
		"No rows": {
			messageThreadID: 2,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT id, is_from_company, is_from_student, content, is_sent, sent_at\s+FROM messages\s+WHERE message_thread_id = \?\s+AND deleted_at IS NULL\s+AND\s+\(\s*\(is_from_company = 1 AND is_sent = 0\)\s+OR \(is_from_company = 1 AND is_sent = 1\)\s+OR \(is_from_student = 1 AND is_sent = 1\)\s*\)\s+ORDER BY sent_at ASC, id ASC;$`).
					WithArgs(int64(2)).
					WillReturnRows(
						sqlmock.NewRows([]string{
							"id", "is_from_company", "is_from_student", "content", "is_sent", "sent_at",
						}),
					)
			},
			wantErr:      false,
			wantMessages: nil,
		},
		"Multiple rows": {
			messageThreadID: 3,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"id", "is_from_company", "is_from_student", "content", "is_sent", "sent_at",
				}).
					AddRow(int64(10), int8(1), int8(0), "Hello", int64(1), time.Date(2025, 1, 1, 12, 0, 0, 0, time.FixedZone("JST", 9*60*60))).
					AddRow(int64(11), int8(0), int8(1), "World", int64(0), time.Date(2025, 1, 1, 12, 5, 0, 0, time.FixedZone("JST", 9*60*60)))

				mock.ExpectQuery(`^SELECT id, is_from_company, is_from_student, content, is_sent, sent_at\s+FROM messages\s+WHERE message_thread_id = \?\s+AND deleted_at IS NULL\s+AND\s+\(\s*\(is_from_company = 1 AND is_sent = 0\)\s+OR \(is_from_company = 1 AND is_sent = 1\)\s+OR \(is_from_student = 1 AND is_sent = 1\)\s*\)\s+ORDER BY sent_at ASC, id ASC;$`).
					WithArgs(int64(3)).
					WillReturnRows(rows)
			},
			wantErr: false,
			wantMessages: entity.Messages{
				&entity.Message{
					ID:            10,
					IsFromCompany: 1,
					IsFromStudent: 0,
					Content:       "Hello",
					IsSent:        1,
					SentAt:        time.Date(2025, 1, 1, 12, 0, 0, 0, time.FixedZone("JST", 9*60*60)),
				},
				&entity.Message{
					ID:            11,
					IsFromCompany: 0,
					IsFromStudent: 1,
					Content:       "World",
					IsSent:        0,
					SentAt:        time.Date(2025, 1, 1, 12, 5, 0, 0, time.FixedZone("JST", 9*60*60)),
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			got, err := mr.GetAllMessagesForCompanyUser(context.Background(), sqlxDB, tc.messageThreadID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantMessages, got)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMessageRepository_GetAllMessagesForStudentUser(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	mr := NewMessageRepository(clock.FixedClocker{})
	tests := map[string]struct {
		messageThreadID entity.MessageThreadID
		mockSetup       func()
		wantErr         bool
		wantMessages    entity.Messages
	}{
		"DB error": {
			messageThreadID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT id, is_from_company, is_from_student, content, is_sent, sent_at\s+FROM messages\s+WHERE message_thread_id = \?\s+AND deleted_at IS NULL\s+AND\s+\(\s*\(is_from_student = 1 AND is_sent = 0\)\s+OR \(is_from_student = 1 AND is_sent = 1\)\s+OR \(is_from_company = 1 AND is_sent = 1\)\s*\)\s+ORDER BY sent_at ASC, id ASC;$`).
					WithArgs(int64(1)).
					WillReturnError(assertAnError())
			},
			wantErr:      true,
			wantMessages: nil,
		},
		"No rows": {
			messageThreadID: 2,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT id, is_from_company, is_from_student, content, is_sent, sent_at\s+FROM messages\s+WHERE message_thread_id = \?\s+AND deleted_at IS NULL\s+AND\s+\(\s*\(is_from_student = 1 AND is_sent = 0\)\s+OR \(is_from_student = 1 AND is_sent = 1\)\s+OR \(is_from_company = 1 AND is_sent = 1\)\s*\)\s+ORDER BY sent_at ASC, id ASC;$`).
					WithArgs(int64(2)).
					WillReturnRows(
						sqlmock.NewRows([]string{
							"id", "is_from_company", "is_from_student", "content", "is_sent", "sent_at",
						}),
					)
			},
			wantErr:      false,
			wantMessages: nil,
		},
		"Multiple rows": {
			messageThreadID: 3,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"id", "is_from_company", "is_from_student", "content", "is_sent", "sent_at",
				}).
					AddRow(int64(10), int8(1), int8(0), "Hello", int64(1), time.Date(2025, 1, 1, 12, 0, 0, 0, time.FixedZone("JST", 9*60*60))).
					AddRow(int64(11), int8(0), int8(1), "World", int64(0), time.Date(2025, 1, 1, 12, 5, 0, 0, time.FixedZone("JST", 9*60*60)))

				mock.ExpectQuery(`^SELECT id, is_from_company, is_from_student, content, is_sent, sent_at\s+FROM messages\s+WHERE message_thread_id = \?\s+AND deleted_at IS NULL\s+AND\s+\(\s*\(is_from_student = 1 AND is_sent = 0\)\s+OR \(is_from_student = 1 AND is_sent = 1\)\s+OR \(is_from_company = 1 AND is_sent = 1\)\s*\)\s+ORDER BY sent_at ASC, id ASC;$`).
					WithArgs(int64(3)).
					WillReturnRows(rows)
			},
			wantErr: false,
			wantMessages: entity.Messages{
				&entity.Message{
					ID:            10,
					IsFromCompany: 1,
					IsFromStudent: 0,
					Content:       "Hello",
					IsSent:        1,
					SentAt:        time.Date(2025, 1, 1, 12, 0, 0, 0, time.FixedZone("JST", 9*60*60)),
				},
				&entity.Message{
					ID:            11,
					IsFromCompany: 0,
					IsFromStudent: 1,
					Content:       "World",
					IsSent:        0,
					SentAt:        time.Date(2025, 1, 1, 12, 5, 0, 0, time.FixedZone("JST", 9*60*60)),
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			got, err := mr.GetAllMessagesForStudentUser(context.Background(), sqlxDB, tc.messageThreadID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantMessages, got)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMessageRepository_AddMessage(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	mr := NewMessageRepository(clock.FixedClocker{})
	tests := map[string]struct {
		inputMessage  *entity.Message
		mockSetup     func(*entity.Message)
		wantErr       bool
		wantID        entity.MessageID
		wantCreatedAt *sql.NullTime
	}{
		"DB error on Exec": {
			inputMessage: &entity.Message{
				MessageThreadID: 100,
				IsFromCompany:   1,
				IsFromStudent:   0,
				Content:         "Hello",
				IsSent:          1,
				SentAt:          time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
				CreatedAt:       clock.FixedClocker{}.Now(),
			},
			mockSetup: func(param *entity.Message) {
				mock.ExpectExec(`^INSERT INTO messages \(message_thread_id, is_from_company, is_from_student, content, is_sent, sent_at, created_at\) VALUES \(\?, \?, \?, \?, \?, \?, \?\);$`).
					WithArgs(
						param.MessageThreadID,
						param.IsFromCompany,
						param.IsFromStudent,
						param.Content,
						param.IsSent,
						param.SentAt,
						param.CreatedAt,
					).
					WillReturnError(assertAnError())
			},
			wantErr: true,
		},
		"LastInsertId error": {
			inputMessage: &entity.Message{
				MessageThreadID: 200,
				IsFromCompany:   0,
				IsFromStudent:   1,
				Content:         "World",
				IsSent:          0,
				SentAt:          time.Date(2025, 2, 1, 10, 30, 0, 0, time.UTC),
				CreatedAt:       clock.FixedClocker{}.Now(),
			},
			mockSetup: func(param *entity.Message) {
				mock.ExpectExec(`^INSERT INTO messages \(message_thread_id, is_from_company, is_from_student, content, is_sent, sent_at, created_at\) VALUES \(\?, \?, \?, \?, \?, \?, \?\);$`).
					WithArgs(
						param.MessageThreadID,
						param.IsFromCompany,
						param.IsFromStudent,
						param.Content,
						param.IsSent,
						param.SentAt,
						param.CreatedAt,
					).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("cannot get lastInsertID")))
			},
			wantErr: true,
		},
		"Success": {
			inputMessage: &entity.Message{
				MessageThreadID: 300,
				IsFromCompany:   1,
				IsFromStudent:   0,
				Content:         "Success case",
				IsSent:          1,
				SentAt:          time.Date(2025, 3, 1, 9, 15, 0, 0, time.UTC),
				CreatedAt:       clock.FixedClocker{}.Now(),
			},
			mockSetup: func(param *entity.Message) {
				mock.ExpectExec(`^INSERT INTO messages \(message_thread_id, is_from_company, is_from_student, content, is_sent, sent_at, created_at\) VALUES \(\?, \?, \?, \?, \?, \?, \?\);$`).
					WithArgs(
						param.MessageThreadID,
						param.IsFromCompany,
						param.IsFromStudent,
						param.Content,
						param.IsSent,
						param.SentAt,
						param.CreatedAt,
					).
					WillReturnResult(sqlmock.NewResult(999, 1))
			},
			wantErr:       false,
			wantID:        999,
			wantCreatedAt: clock.FixedClocker{}.Now(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup(tc.inputMessage)
			err := mr.AddMessage(context.Background(), sqlxDB, tc.inputMessage)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantID, tc.inputMessage.ID)
				assert.Equal(t, tc.wantCreatedAt, tc.inputMessage.CreatedAt)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMessageRepository_EditMessage(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	mr := NewMessageRepository(clock.FixedClocker{})
	tests := map[string]struct {
		inputMessage  *entity.Message
		mockSetup     func(*entity.Message)
		wantErr       bool
		wantUpdatedAt *sql.NullTime
	}{
		"DB error on Exec": {
			inputMessage: &entity.Message{
				ID:        1,
				Content:   "Before Update - DB error",
				UpdatedAt: clock.FixedClocker{}.Now(),
			},
			mockSetup: func(m *entity.Message) {
				mock.ExpectExec(`^UPDATE messages SET content = \?, updated_at = \? WHERE id = \?;$`).
					WithArgs(
						m.Content,
						m.UpdatedAt,
						m.ID,
					).
					WillReturnError(assertAnError())
			},
			wantErr: true,
		},
		"Success": {
			inputMessage: &entity.Message{
				ID:        2,
				Content:   "Update Content",
				UpdatedAt: clock.FixedClocker{}.Now(),
			},
			mockSetup: func(m *entity.Message) {
				mock.ExpectExec(`^UPDATE messages SET content = \?, updated_at = \? WHERE id = \?;$`).
					WithArgs(
						m.Content,
						m.UpdatedAt,
						m.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1)) // InsertID=0, RowsAffected=1
			},
			wantErr:       false,
			wantUpdatedAt: clock.FixedClocker{}.Now(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup(tc.inputMessage)
			err := mr.EditMessage(context.Background(), sqlxDB, tc.inputMessage)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantUpdatedAt, tc.inputMessage.UpdatedAt)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMessageRepository_DeleteMessage(t *testing.T) {
	sqlxDB, mock := newMockDB(t)
	mr := NewMessageRepository(clock.FixedClocker{})
	tests := map[string]struct {
		messageID entity.MessageID
		mockSetup func(entity.MessageID)
		wantErr   bool
	}{
		"DB error": {
			messageID: 1,
			mockSetup: func(id entity.MessageID) {
				mock.ExpectExec(`^UPDATE messages SET deleted_at = \? WHERE id = \?;$`).
					WithArgs(clock.FixedClocker{}.Now(), id).
					WillReturnError(assertAnError())
			},
			wantErr: true,
		},
		"Success": {
			messageID: 2,
			mockSetup: func(id entity.MessageID) {
				mock.ExpectExec(`^UPDATE messages SET deleted_at = \? WHERE id = \?;$`).
					WithArgs(clock.FixedClocker{}.Now(), id).
					WillReturnResult(sqlmock.NewResult(0, 1)) // InsertID=0, RowsAffected=1
			},
			wantErr: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup(tc.messageID)
			err := mr.DeleteMessage(context.Background(), sqlxDB, tc.messageID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
