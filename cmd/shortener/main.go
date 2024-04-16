package main

import (
	"krajcik/shortener/cmd/shortener/handler"
	"krajcik/shortener/internal/app/shortener"
	"net/http"
)

var service = shortener.NewService(shortener.NewRepository())

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(":8080", http.HandlerFunc(webhook))
}

func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handler.PostShrt(service)(w, r)
		return
	}

	handler.GetShrt(service)(w, r)
}
