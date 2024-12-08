package credential

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GenerateAccessToken(appKind string, userID int64) (string, error) {
	secretKey, err := getAccessTokenSecretKey()
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	randomStr := generateRandomStr(16)
	baseData := fmt.Sprintf("appkind:%s|user_id:%d|random:%s", appKind, userID, randomStr)
	cipherText := aesGCM.Seal(nil, nonce, []byte(baseData), nil)
	combined := append(nonce, cipherText...)
	return base64.StdEncoding.EncodeToString(combined), nil
}

func getAccessTokenSecretKey() ([]byte, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	secretKeyBase64 := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	if secretKeyBase64 == "" {
		return nil, fmt.Errorf("ACCESS_TOKEN_SECRET_KEY environment variable not set")
	}
	return base64.StdEncoding.DecodeString(secretKeyBase64)
}
