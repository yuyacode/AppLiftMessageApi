package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"

	"github.com/yuyacode/AppLiftMessageApi/clock"
	"github.com/yuyacode/AppLiftMessageApi/config"
	"github.com/yuyacode/AppLiftMessageApi/handler"
	"github.com/yuyacode/AppLiftMessageApi/service"
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
	clocker := clock.RealClocker{}
	messageRepo := &store.MessageRepository{
		Clocker: clocker,
	}
	oAuthRepo := &store.OAuthRepository{
		Clocker: clocker,
	}

	ro := &handler.RegisterOAuth{
		Service: &service.RegisterOAuth{
			DBHandlers:       dbHandlerList,
			CredentialGetter: oAuthRepo,
			CredentialSetter: oAuthRepo,
		},
		Validator: v,
	}
	mux.Route("/oauth", func(r chi.Router) {
		r.Get("/register", ro.ServeHTTP)
	})

	gm := &handler.GetMessage{
		Service: &service.GetMessage{
			DBHandlers:         dbHandlerList,
			MessageGetter:      messageRepo,
			MessageOwnerGetter: messageRepo,
		},
		Validator: v,
	}
	mux.Route("/", func(r chi.Router) {
		// ミドルウェア通す
		r.Get("/", gm.ServeHTTP)
	})

	return mux, dbCloseFuncList, nil
}
