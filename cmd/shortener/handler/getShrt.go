package handler

import (
	"errors"
	"krajcik/shortener/internal/app/shortener"
	"net/http"
	"strings"
)

func GetShrt(s *shortener.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := strings.Trim(r.URL.Path, "/")
		shrt, err := s.URLByShrt(url)

		if err != nil {
			if errors.Is(err, shortener.ErrNotFound) {
				http.NotFound(w, r)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", shrt)
		w.WriteHeader(http.StatusTemporaryRedirect)
		_, _ = w.Write([]byte(shrt))
	}
}
