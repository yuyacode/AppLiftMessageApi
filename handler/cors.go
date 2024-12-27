package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func CORSMiddleware() func(http.Handler) http.Handler {
	allowedOrigin, err := getAllowedOrigin()
	if err != nil {
		log.Printf("failed to fetch allowed origin: %v", err)
		allowedOrigin = "*" // フォールバックとして全オリジンを許可
	}
	return cors.Handler(cors.Options{
		AllowedOrigins: []string{allowedOrigin},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:         300,
	})
}

func getAllowedOrigin() (string, error) {
	if err := godotenv.Load(); err != nil {
		if os.Getenv("ENV") == "dev" {
			return "", err
		}
	}
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		return "", fmt.Errorf("ALLOWED_ORIGIN environment variable not set")
	}
	return allowedOrigin, nil
}
