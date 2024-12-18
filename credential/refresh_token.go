package credential

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/yuyacode/AppLiftMessageApi/handler"
)

func GenerateRefreshToken(appKind string, userID int64) (string, error) {
	secretKey, err := getRefreshTokenSecretKey()
	if err != nil {
		return "", fmt.Errorf("failed to get refresh token secret key: %w", err)
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
	randomStr := generateRandomStr(32)
	baseData := fmt.Sprintf("appkind:%s|user_id:%d|random:%s", appKind, userID, randomStr)
	cipherText := aesGCM.Seal(nil, nonce, []byte(baseData), nil)
	combined := append(nonce, cipherText...)
	return base64.StdEncoding.EncodeToString(combined), nil
}

func DecryptRefreshToken(refreshToken string) (string, int64, error) {
	decoded, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		return "", 0, handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to decode refresh token",
			err.Error(),
		)
	}
	secretKey, err := getRefreshTokenSecretKey()
	if err != nil {
		return "", 0, handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to get refresh token secret key",
			err.Error(),
		)
	}
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", 0, handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to create cipher block",
			err.Error(),
		)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", 0, handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to create GCM",
			err.Error(),
		)
	}
	if len(decoded) < aesGCM.NonceSize() {
		return "", 0, handler.NewServiceError(
			http.StatusUnauthorized,
			"invalid_token",
			"refresh token too short",
		)
	}
	nonce, cipherText := decoded[:aesGCM.NonceSize()], decoded[aesGCM.NonceSize():]
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", 0, handler.NewServiceError(
			http.StatusInternalServerError,
			"failed to decrypt refresh token",
			err.Error(),
		)
	}
	parsedData := string(plainText)
	parts := strings.Split(parsedData, "|")
	if len(parts) != 3 {
		return "", 0, handler.NewServiceError(
			http.StatusUnauthorized,
			"invalid_token",
			"invalid refresh token format",
		)
	}
	var appKind string
	if strings.HasPrefix(parts[0], "appkind:") {
		appKind = strings.TrimPrefix(parts[0], "appkind:")
		if appKind != "company" && appKind != "student" {
			return "", 0, handler.NewServiceError(
				http.StatusUnauthorized,
				"invalid_token",
				fmt.Sprintf("invalid appkind in refresh token: %v", err),
			)
		}
	} else {
		return "", 0, handler.NewServiceError(
			http.StatusUnauthorized,
			"invalid_token",
			"appkind not found in refresh token",
		)
	}
	var userID int64
	if strings.HasPrefix(parts[1], "user_id:") {
		userIDStr := strings.TrimPrefix(parts[1], "user_id:")
		userID, err = strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return "", 0, handler.NewServiceError(
				http.StatusUnauthorized,
				"invalid_token",
				fmt.Sprintf("invalid user_id in refresh token: %v", err),
			)
		}
	} else {
		return "", 0, handler.NewServiceError(
			http.StatusUnauthorized,
			"invalid_token",
			"user_id not found in refresh token",
		)
	}
	return appKind, userID, nil
}

func getRefreshTokenSecretKey() ([]byte, error) {
	if err := godotenv.Load("../.env"); err != nil {
		if os.Getenv("ENV") == "dev" {
			return nil, err
		}
	}
	secretKeyBase64 := os.Getenv("REFRESH_TOKEN_SECRET_KEY")
	if secretKeyBase64 == "" {
		return nil, fmt.Errorf("REFRESH_TOKEN_SECRET_KEY environment variable not set")
	}
	secretKey, err := base64.StdEncoding.DecodeString(secretKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode refresh token secret key: %w", err)
	}
	return secretKey, nil
}
