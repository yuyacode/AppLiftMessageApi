//go:build batch

package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/yuyacode/AppLiftMessageApi/batch"
)

func init() {
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("failed to set time.Local: %v", err)
	}
}

// example: go run -tags=batch batch.go --mode=generate_api_key --target=company
func main() {
	mode := flag.String("mode", "", "mode: 'generate_api_key' or 'generate_access_token_secret_key' or 'generate_refresh_token_secret_key'")
	target := flag.String("target", "", "target: 'company' or 'student'")
	flag.Parse()
	switch *mode {
	case "generate_api_key":
		if *target == "" {
			log.Fatalf("missing required option '--target'")
		}
		if *target != "company" && *target != "student" {
			log.Fatalf("invalid target")
		}
		apiKey, err := batch.GenerateAPIKey(*target)
		if err != nil {
			log.Fatalf("failed to generate API key: %v", err)
		}
		fmt.Printf("API Key successfully generated: %s\n", apiKey)
	case "generate_access_token_secret_key", "generate_refresh_token_secret_key":
		if *target != "" {
			log.Fatalf("unnecessary option '--target'")
		}
		secretKey, err := batch.GenerateTokenSecretKey()
		if err != nil {
			log.Fatalf("failed to %s: %v", *mode, err)
		}
		fmt.Printf("succeeded to %s: %s\n", *mode, secretKey)
	default:
		log.Fatalf("invalid mode")
	}
}
