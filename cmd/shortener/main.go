package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"krajcik/shortener/cmd/shortener/config"
	"krajcik/shortener/cmd/shortener/handler"
	"krajcik/shortener/internal/app/shortener"
	internalmiddleware "krajcik/shortener/internal/middleware"
	"net/http"
	"time"
)

var params *config.Params

var service *shortener.Service
var logger *zap.Logger

func main() {
	if err := run(); err != nil {
		logger.Panic(err.Error())
	}
}

func run() error {
	params = config.Create()
	return http.ListenAndServe(params.A, router())
}

func init() {
	service = shortener.NewService(shortener.NewRepository())
	err := initLogger("debug")
	if err != nil {
		panic(err)
	}
}

func router() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(internalmiddleware.Logger(logger, ""))
	r.Use(middleware.NoCache)
	r.Use(middleware.SetHeader("Content-Type", "plain/text; charset=utf-8"))
	r.Use(middleware.Recoverer)

	r.Get("/{shrt}", handler.GetShrt(service))
	r.Post("/", handler.PostShrt(service, params))

	return r
}

func initLogger(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	logger, err = cfg.Build()
	if err != nil {
		return err
	}

	return nil
}
