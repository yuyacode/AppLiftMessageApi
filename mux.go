package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"

	"github.com/yuyacode/AppLiftMessageApi/config"
	"github.com/yuyacode/AppLiftMessageApi/store"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, map[string]func(), error) {
	mux := chi.NewRouter()
	v := validator.New()
	targetDBList := [3]string{"company", "student", "common"}
	var dbHandlerList = make(map[string]*sqlx.DB, len(targetDBList))
	var dbCloseFuncList = make(map[string]func(), len(targetDBList))
	var err error
	for _, v := range targetDBList {
		dbHandlerList[v], dbCloseFuncList[v], err = store.New(ctx, cfg, v)
		if err != nil {
			return nil, dbCloseFuncList, err
		}
	}
	// muxの定義
	return mux, dbCloseFuncList, nil
}
