//go:build testdata

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/yuyacode/AppLiftMessageApi/credential"
)

// example: go run -tags=testdata testdata.go --action=??? arg1 arg2 arg3
func main() {
	action := flag.String("action", "", "Specifies the function or method to execute.")
	flag.Parse()
	args := flag.Args()
	switch *action {
	case "hash_api_key":
		if len(args) < 1 {
			log.Fatalf("hash_api_key requires an argument [apiKey]. Usage:\n  go run -tags=testdata testdata.go --action=hash_api_key your_api_key\n")
		}
		fmt.Printf("API Key successfully hashed: %s\n", credential.HashAPIKey(args[0]))
	default:
		log.Fatalf("invalid action")
	}
}
