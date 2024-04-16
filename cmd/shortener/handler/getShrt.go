package handler

import (
	"errors"
	"krajcik/shortener/internal/app/shortener"
	"net/http"
	"strings"
)

func GetShrt(s *shortener.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		url := strings.Trim(r.URL.Path, "/")
		shrt, err := s.UrlByShrt(url)

		if err != nil {
			if errors.Is(err, shortener.NotFoundError) {
				http.NotFound(w, r)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(shrt))
	}
}
