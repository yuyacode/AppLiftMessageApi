package entity

import (
	"time"
)

type MessageThreadID int64

type MessageThread struct {
	ID            MessageThreadID `json:"id"              db:"id"`
	CompanyUserID int64           `json:"company_user_id" db:"company_user_id"`
	StudentUserID int64           `json:"student_user_id" db:"student_user_id"`
	CreatedAt     time.Time       `json:"created_at"      db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"      db:"updated_at"`
	DeletedAt     time.Time       `json:"deleted_at"      db:"deleted_at"`
}
