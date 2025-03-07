package clock

import (
	"database/sql"
	"time"
)

type Clocker interface {
	Now() *sql.NullTime
}

type RealClocker struct{}

func (rc RealClocker) Now() *sql.NullTime {
	return &sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
}

type FixedClocker struct{}

func (fc FixedClocker) Now() *sql.NullTime {
	jst := time.FixedZone("JST", 9*60*60)
	return &sql.NullTime{
		Time:  time.Date(2025, 1, 1, 9, 0, 0, 0, jst),
		Valid: true,
	}
}
