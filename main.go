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
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	mux, dbCloseFuncs, err := NewMux(ctx, cfg)
	if err != nil {
		for _, f := range dbCloseFuncs {
			f()
		}
		return err
	}
	for _, f := range dbCloseFuncs {
		defer func(f func()) {
			f()
		}(f)
	}
	// サーバーを生成して起動する
}
