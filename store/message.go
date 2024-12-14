package store

import (
	"context"

	"github.com/yuyacode/AppLiftMessageApi/clock"
	"github.com/yuyacode/AppLiftMessageApi/entity"
)

type MessageRepository struct {
	Clocker clock.Clocker
}

func NewMessageRepository(clocker clock.Clocker) *MessageRepository {
	return &MessageRepository{
		Clocker: clocker,
	}
}

func (mr *MessageRepository) GetAllMessages(ctx context.Context, db Queryer, threadID entity.MessageThreadID) (entity.Messages, error) {
	query := `
        SELECT id, message_thread_id, is_from_company, is_from_student, content, is_unread, created_at, updated_at, deleted_at
        FROM messages
        WHERE message_thread_id = ? AND deleted_at IS NULL
		ORDER BY id ASC;
    `
	rows, err := db.QueryxContext(ctx, query, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages entity.Messages
	for rows.Next() {
		var m entity.Message
		if err := rows.StructScan(&m); err != nil {
			return nil, err
		}
		messages = append(messages, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (mr *MessageRepository) GetThreadCompanyOwner(ctx context.Context, db Queryer, param *entity.MessageThread) (int64, error) {
	query := "SELECT company_user_id FROM message_threads WHERE id = :id AND deleted_at IS NULL;"
	var companyUserID int64
	if err := db.GetContext(ctx, &companyUserID, query, param); err != nil {
		return 0, err
	}
	return companyUserID, nil
}
