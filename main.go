package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/yuyacode/AppLiftMessageApi/config"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	mux, dbCloseFuncList, err := NewMux(ctx, cfg)
	if err != nil {
		for _, dbCloseFunc := range dbCloseFuncList {
			dbCloseFunc()
		}
		return err
	}
	for _, dbCloseFunc := range dbCloseFuncList {
		defer func(f func()) {
			f()
		}(dbCloseFunc)
	}
	// サーバーを生成して起動する
}
