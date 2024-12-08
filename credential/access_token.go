package credential

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func GenerateAccessToken(appKind string, userID int64) (string, error) {
	secretKey, err := getAccessTokenSecretKey()
	if err != nil {
		return "", fmt.Errorf("failed to get access token secret key: %w", err)
	}
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %w", err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	randomStr := generateRandomStr(16)
	baseData := fmt.Sprintf("appkind:%s|user_id:%d|random:%s", appKind, userID, randomStr)
	cipherText := aesGCM.Seal(nil, nonce, []byte(baseData), nil)
	combined := append(nonce, cipherText...)
	return base64.StdEncoding.EncodeToString(combined), nil
}

func DecryptAccessToken(accessToken string) (string, int64, error) {
	decoded, err := base64.StdEncoding.DecodeString(accessToken)
	if err != nil {
		return "", 0, fmt.Errorf("failed to decode access token: %w", err)
	}
	secretKey, err := getAccessTokenSecretKey()
	if err != nil {
		return "", 0, fmt.Errorf("failed to get access token secret key: %w", err)
	}
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create cipher block: %w", err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create GCM: %w", err)
	}
	if len(decoded) < aesGCM.NonceSize() {
		return "", 0, errors.New("access token too short")
	}
	nonce, cipherText := decoded[:aesGCM.NonceSize()], decoded[aesGCM.NonceSize():]
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", 0, fmt.Errorf("failed to decrypt access token: %w", err)
	}
	parsedData := string(plainText)
	parts := strings.Split(parsedData, "|")
	if len(parts) != 3 {
		return "", 0, errors.New("invalid access token format")
	}
	var appKind string
	if strings.HasPrefix(parts[0], "appkind:") {
		appKind = strings.TrimPrefix(parts[0], "appkind:")
	} else {
		return "", 0, errors.New("appkind not found in access token")
	}
	var userID int64
	if strings.HasPrefix(parts[1], "user_id:") {
		userIDStr := strings.TrimPrefix(parts[1], "user_id:")
		userID, err = strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return "", 0, fmt.Errorf("invalid user_id in access token: %w", err)
		}
	} else {
		return "", 0, errors.New("user_id not found in access token")
	}
	return appKind, userID, nil
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
