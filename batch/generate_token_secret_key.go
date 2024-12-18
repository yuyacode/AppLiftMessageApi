package batch

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateTokenSecretKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	encodedKey := base64.StdEncoding.EncodeToString(key)
	return encodedKey, nil
}
