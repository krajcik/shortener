package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"krajcik/shortener/cmd/shortener/config"
	"krajcik/shortener/cmd/shortener/handler"
	"krajcik/shortener/cmd/shortener/handler/api"
	"krajcik/shortener/internal/app/shortener"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	gzm "krajcik/shortener/cmd/shortener/middleware/gzip"
	internalmiddleware "krajcik/shortener/internal/middleware"
)

var service *shortener.Service
var logger *zap.Logger

func main() {
	if err := run(); err != nil {
		panic(err)
		//logger.Panic(err.Error())
	}
}

func run() error {
	params, err := config.Create()
	db, err := db(params)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	r, err := router(params, db)
	if err != nil {
		return err
	}
	return http.ListenAndServe(params.A, r)
}

func router(params *config.Params, db *sql.DB) (chi.Router, error) {
	repository := shortener.NewDbRepository(db)

	service = shortener.NewService(repository)
	if err := initLogger("debug"); err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Heartbeat("/p"))
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(internalmiddleware.Logger(logger, ""))
	r.Use(middleware.NoCache)
	r.Use(middleware.SetHeader("Content-Type", "text/plain; charset=utf-8"))
	r.Use(gzm.Middleware(logger))
	r.Use(middleware.Recoverer)
	//r.Use(middleware.Compress())

	r.Get("/{shrt}", handler.GetShrt(service))
	r.Post("/", handler.PostShrt(service, params))
	r.Get("/ping", handler.Ping(db, logger))

	r.Mount("/api", apiRouter(params))

	return r, nil
}

func db(params *config.Params) (*sql.DB, error) {
	db, err := sql.Open("pgx", params.DatabaseDsn)
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}

func apiRouter(params *config.Params) http.Handler {
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
