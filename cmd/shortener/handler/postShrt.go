package handler

import (
	"fmt"
	"io"
	"krajcik/shortener/internal/app/shortener"
	"net/http"
)

func PostShrt(s *shortener.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		res, err := s.ShrtByURL(string(buf))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res = fmt.Sprintf("%s/%s", r.Host, res)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(res))
		if err != nil {
			panic(err)
		}
	}
}
