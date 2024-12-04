package entity

import (
	"time"
)

type MessageAPICredentialID int64

type MessageAPICredential struct {
	ID           MessageAPICredentialID `json:"id"            db:"id"`
	UserID       int64                  `json:"user_id"       db:"user_id"`
	ClientID     string                 `json:"client_id"     db:"client_id"`
	ClientSecret string                 `json:"client_secret" db:"client_secret"`
	AccessToken  string                 `json:"access_token"  db:"access_token"`
	RefreshToken string                 `json:"refresh_token" db:"refresh_token"`
	ExpiresAt    time.Time              `json:"expires_at"    db:"expires_at"`
	CreatedAt    time.Time              `json:"created_at"    db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"    db:"updated_at"`
	DeletedAt    time.Time              `json:"deleted_at"    db:"deleted_at"`
}
