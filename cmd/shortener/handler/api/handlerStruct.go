package api

import (
	"krajcik/shortener/cmd/shortener/config"
	"krajcik/shortener/internal/app/shortener"

	"go.uber.org/zap" //todo switch to logger interface
)

type PostShrtHandler struct {
	S *shortener.Service
	P *config.Params
	L *zap.Logger
}

type ShortenBatchHandler struct {
	S *shortener.Service
	P *config.Params
	L *zap.Logger
}
