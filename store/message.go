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

func (mr *MessageRepository) GetThreadCompanyOwner(ctx context.Context, db Queryer, messageThreadID entity.MessageThreadID) (int64, error) {
	query := "SELECT company_user_id FROM message_threads WHERE id = ? AND deleted_at IS NULL;"
	var companyUserID int64
	if err := db.GetContext(ctx, &companyUserID, query, messageThreadID); err != nil {
		return 0, err
	}
	return companyUserID, nil
}

func (mr *MessageRepository) GetThreadStudentOwner(ctx context.Context, db Queryer, messageThreadID entity.MessageThreadID) (int64, error) {
	query := "SELECT student_user_id FROM message_threads WHERE id = ? AND deleted_at IS NULL;"
	var studentUserID int64
	if err := db.GetContext(ctx, &studentUserID, query, messageThreadID); err != nil {
		return 0, err
	}
	return studentUserID, nil
}

func (mr *MessageRepository) GetThreadCompanyOwnerByMessageID(ctx context.Context, db Queryer, messageID entity.MessageID) (int64, error) {
	query := `
		SELECT company_user_id
		FROM message_threads
		INNER JOIN messages
		ON message_threads.id = messages.message_thread_id
		WHERE messages.id = ? AND is_from_company = 1;
	`
	var companyUserID int64
	if err := db.GetContext(ctx, &companyUserID, query, messageID); err != nil {
		return 0, err
	}
	return companyUserID, nil
}

func (mr *MessageRepository) GetAllMessagesForCompanyUser(ctx context.Context, db Queryer, messageThreadID entity.MessageThreadID) (entity.Messages, error) {
	query := `
        SELECT id, is_from_company, is_from_student, content, is_sent, sent_at
        FROM messages
        WHERE message_thread_id = ?
		AND deleted_at IS NULL
		AND
		(
			(is_from_company = 1 AND is_sent = 0)
      		OR (is_from_company = 1 AND is_sent = 1)
      		OR (is_from_student = 1 AND is_sent = 1)
		)
		ORDER BY sent_at ASC, id ASC;
    `
	rows, err := db.QueryxContext(ctx, query, messageThreadID)
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

func (mr *MessageRepository) GetAllMessagesForStudentUser(ctx context.Context, db Queryer, messageThreadID entity.MessageThreadID) (entity.Messages, error) {
	query := `
        SELECT id, is_from_company, is_from_student, content, is_sent, sent_at
        FROM messages
        WHERE message_thread_id = ?
		AND deleted_at IS NULL
		AND
		(
			(is_from_student = 1 AND is_sent = 0)
      		OR (is_from_student = 1 AND is_sent = 1)
      		OR (is_from_company = 1 AND is_sent = 1)
		)
		ORDER BY sent_at ASC, id ASC;
    `
	rows, err := db.QueryxContext(ctx, query, messageThreadID)
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

func (mr *MessageRepository) AddMessage(ctx context.Context, db Execer, param *entity.Message) error {
	param.CreatedAt = mr.Clocker.Now()
	query := "INSERT INTO messages (message_thread_id, is_from_company, is_from_student, content, is_sent, sent_at, created_at) VALUES (:message_thread_id, :is_from_company, :is_from_student, :content, :is_sent, :sent_at, :created_at);"
	result, err := db.NamedExecContext(ctx, query, param)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	param.ID = entity.MessageID(id)
	return nil
}

func (mr *MessageRepository) EditMessage(ctx context.Context, db Execer, param *entity.Message) error {
	param.UpdatedAt = mr.Clocker.Now()
	query := "UPDATE messages SET content = :content, updated_at = :updated_at WHERE id = :id;"
	_, err := db.NamedExecContext(ctx, query, param)
	if err != nil {
		return err
	}
	return nil
}

func (mr *MessageRepository) DeleteMessage(ctx context.Context, db Execer, id entity.MessageID) error {
	query := "UPDATE messages SET deleted_at = ? WHERE id = ?;"
	_, err := db.ExecContext(ctx, query, mr.Clocker.Now(), id)
	if err != nil {
		return err
	}
	return nil
}
