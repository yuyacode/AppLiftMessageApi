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
	dbList := [3]string{"company", "student", "common"}
	var dbHandlers = make(map[string]*sqlx.DB, len(dbList))
	var dbCloseFuncs = make(map[string]func(), len(dbList))
	var err error
	for _, v := range dbList {
		dbHandlers[v], dbCloseFuncs[v], err = store.New(ctx, cfg, v)
		if err != nil {
			return nil, dbCloseFuncs, err
		}
	}
	clocker := clock.RealClocker{}
	oAuthRepo := store.NewOAuthRepository(clocker)
	roService := service.NewRegisterOAuth(dbHandlers, oAuthRepo, oAuthRepo)
	roHandler := handler.NewRegisterOAuth(roService, v)
	mux.Route("/messages", func(r chi.Router) {
		r.Post("/register", roHandler.ServeHTTP)
		// r.With(RefreshTokenMiddleware).Post("/token", fooHandler.ServeHTTP)
	})
	messageRepo := store.NewMessageRepository(clocker)
	gmService := service.NewGetMessage(dbHandlers, messageRepo, messageRepo)
	gmHandler := handler.NewGetMessage(gmService, v)
	mux.Route("/messages", func(r chi.Router) {
		// r.Use(AccessTokenMiddleware)
		r.Get("/", gmHandler.ServeHTTP)
		// r.Post("/", fooHandler.ServeHTTP)
		// r.Patch("/{id}", fooHandler.ServeHTTP)
		// r.Delete("/{id}", fooHandler.ServeHTTP)
	})
	return mux, dbCloseFuncs, nil
}
