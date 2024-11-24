package entity

import (
	"time"
)

type MessageID int64

type Message struct {
	ID              MessageID       `json:"id"                db:"id"`
	MessageThreadID MessageThreadID `json:"message_thread_id" db:"message_thread_id"`
	IsFromCompany   int8            `json:"is_from_company"   db:"is_from_company"`
	IsFromStudent   int8            `json:"is_from_student"   db:"is_from_student"`
	Content         string          `json:"content"           db:"content"`
	IsUnread        string          `json:"is_unread"         db:"is_unread"`
	CreatedAt       time.Time       `json:"created_at"        db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"        db:"updated_at"`
	DeletedAt       time.Time       `json:"deleted_at"        db:"deleted_at"`
}

type Messages []*Message
