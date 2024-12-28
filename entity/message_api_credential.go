package entity

import (
	"database/sql"
)

type MessageAPICredentialID int64

type MessageAPICredential struct {
	ID           MessageAPICredentialID `json:"id"            db:"id"`
	UserID       int64                  `json:"user_id"       db:"user_id"`
	ClientID     string                 `json:"client_id"     db:"client_id"`
	ClientSecret string                 `json:"client_secret" db:"client_secret"`
	AccessToken  string                 `json:"access_token"  db:"access_token"`
	RefreshToken string                 `json:"refresh_token" db:"refresh_token"`
	ExpiresAt    *sql.NullTime          `json:"expires_at"    db:"expires_at"`
	CreatedAt    *sql.NullTime          `json:"created_at"    db:"created_at"`
	UpdatedAt    *sql.NullTime          `json:"updated_at"    db:"updated_at"`
	DeletedAt    *sql.NullTime          `json:"deleted_at"    db:"deleted_at"`
}
