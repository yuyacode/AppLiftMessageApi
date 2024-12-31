package entity

import (
	"database/sql"
	"time"
)

type MessageID int64

type Message struct {
	ID              MessageID       `json:"id"                db:"id"`
	MessageThreadID MessageThreadID `json:"message_thread_id" db:"message_thread_id"`
	IsFromCompany   int8            `json:"is_from_company"   db:"is_from_company"`
	IsFromStudent   int8            `json:"is_from_student"   db:"is_from_student"`
	Content         string          `json:"content"           db:"content"`
	IsSent          int8            `json:"is_sent"           db:"is_sent"`
	SentAt          time.Time       `json:"sent_at"           db:"sent_at"`
	CreatedAt       *sql.NullTime   `json:"created_at"        db:"created_at"`
	UpdatedAt       *sql.NullTime   `json:"updated_at"        db:"updated_at"`
	DeletedAt       *sql.NullTime   `json:"deleted_at"        db:"deleted_at"`
}

type Messages []*Message
