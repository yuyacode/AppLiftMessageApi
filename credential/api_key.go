package credential

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashAPIKey(apiKey string) string {
	hasher := sha256.New()
	hasher.Write([]byte(apiKey))
	return hex.EncodeToString(hasher.Sum(nil))
}
