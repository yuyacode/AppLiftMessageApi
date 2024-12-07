package clock

import (
	"time"
)

type Clocker interface {
	Now() time.Time
}

type RealClocker struct{}

func (r RealClocker) Now() time.Time {
	return time.Now()
}

type FixedClocker struct{}

func (fc FixedClocker) Now() time.Time {
	jst := time.FixedZone("JST", 9*60*60)
	return time.Date(2024, 7, 20, 12, 33, 56, 0, jst)
}
