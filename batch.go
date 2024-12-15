//go:build batch

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/yuyacode/AppLiftMessageApi/batch"
)

// example: go run -tags=batch batch.go --mode=generate_api_key --target=company
func main() {
	mode := flag.String("mode", "", "mode: 'generate_api_key' or 'other...'")
	target := flag.String("target", "", "target: 'company' or 'student'")
	flag.Parse()
	if *mode == "" || *target == "" {
		log.Fatalf("missing required option '--mode' or '--target'")
	}
	if *mode != "generate_api_key" {
		log.Fatalf("invalid mode")
	}
	if *target != "company" && *target != "student" {
		log.Fatalf("invalid target")
	}
	switch *mode {
	case "generate_api_key":
		apiKey, err := batch.GenerateAPIKey(*target)
		if err != nil {
			log.Fatalf("failed to generate API key: %v", err)
		}
		fmt.Printf("API Key successfully generated: %s\n", apiKey)
	}
}
