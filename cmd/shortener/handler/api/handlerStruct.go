package api

import (
	"go.uber.org/zap" //todo switch to logger interface
	"krajcik/shortener/cmd/shortener/config"
	"krajcik/shortener/internal/app/shortener"
)

type PostShrtHandler struct {
	S *shortener.Service
	P *config.Params
	L *zap.Logger
}
