package main

import (
	"net/http"
	"time"

	"krajcik/shortener/cmd/shortener/config"
	"krajcik/shortener/cmd/shortener/handler"
	"krajcik/shortener/cmd/shortener/handler/api"
	"krajcik/shortener/internal/app/shortener"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	gzm "krajcik/shortener/cmd/shortener/middleware/gzip"
	internalmiddleware "krajcik/shortener/internal/middleware"
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
	r.Use(middleware.SetHeader("Content-Type", "text/plain; charset=utf-8"))
	r.Use(gzm.Middleware(logger))
	r.Use(middleware.Recoverer)
	//r.Use(middleware.Compress())

	r.Get("/{shrt}", handler.GetShrt(service))
	r.Post("/", handler.PostShrt(service, params))

	r.Mount("/api", apiRouter())

	return r
}

func apiRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.SetHeader("Content-Type", "application/json; charset=utf-8"))
	r.Use(middleware.AllowContentType("application/json"))
	shrtHandler := &api.PostShrtHandler{
		S: service,
		P: params,
		L: logger,
	}
	r.Post("/shorten", http.HandlerFunc(shrtHandler.PostShrt))

	//r.Route("/{articleID}", func(r chi.Router) {
	//	r.Get("/", getArticle)
	// r.Put("/", updateArticle)
	// r.Delete("/", deleteArticle)
	//})
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
