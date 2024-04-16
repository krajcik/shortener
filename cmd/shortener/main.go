package main

import (
	"errors"
	"io"
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
		if r.Header.Get("Content-Type") != "text/plain" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body := r.Body
		defer func(body io.ReadCloser) {
			err := body.Close()
			if err != nil {
				panic(err)
			}
		}(body)

		buf, err := io.ReadAll(body)
		if err != nil {
			panic(err)
		}

		res, err := service.ShrtByUrl(string(buf))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//w.Header().Set("Content-Type", "text/plain")
		_, err = w.Write([]byte(res))
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusCreated)
		return
	}

	// установим правильный заголовок для типа данных
	w.Header().Set("Content-Type", "text/plain")

	shrt, err := service.UrlByShrt(r.URL.RawQuery)
	if err != nil {
		if errors.Is(err, shortener.NotFoundError) {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte(shrt))
	w.WriteHeader(http.StatusCreated)
}
