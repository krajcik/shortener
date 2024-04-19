package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"krajcik/shortener/cmd/shortener/config"
	"krajcik/shortener/cmd/shortener/handler"
	"krajcik/shortener/internal/app/shortener"
	"net/http"
	"time"
)

var params *config.Params

var service = shortener.NewService(shortener.NewRepository())

func main() {
	params = config.Create()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	return http.ListenAndServe(params.A, router())
}

func router() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.NoCache)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.SetHeader("Content-Type", "plain/text; charset=utf-8"))
	r.Use(middleware.Recoverer)

	r.Get("/{shrt}", handler.GetShrt(service))
	r.Post("/", handler.PostShrt(service, params))

	return r
}
