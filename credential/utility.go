package credential

import (
	"math/rand"
	"strings"
	"time"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func generateRandomStr(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	builder := strings.Builder{}
	builder.Grow(length)
	for i := 0; i < length; i++ {
		randomIndex := seededRand.Intn(len(charset))
		builder.WriteByte(charset[randomIndex])
	}
	return builder.String()
}
