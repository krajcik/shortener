package gzip

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

var allowedTypes = []string{
	"text/html",
	"application/json",
}

func Middleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ow := w
				isAllowedContent := false
				for _, v := range allowedTypes {
					if strings.Contains(r.Header.Get("Content-Type"), v) {
						isAllowedContent = true
						break
					}
				}
				isGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
				logger.Debug(
					":",
					zap.Bool("isAllowedContent", isAllowedContent),
					zap.Bool("isGzip", isGzip),
				)
				if isAllowedContent && isGzip {
					logger.Debug("Decompressing request")
					w.Header().Set("Content-Encoding", "gzip")
					gz := newGzipResponseWriter(w)
					defer func(gz *gzipResponseWriter) {
						err := gz.Close()
						if err != nil {
							logger.Error("Failed to close gzip writer", zap.Error(err))
						}
					}(gz)
					ow = gz
				}

				contentEncoding := r.Header.Get("Content-Encoding")
				sendsGzip := strings.Contains(contentEncoding, "gzip")
				if sendsGzip {
					gr, err := newGzipReader(r.Body)
					if err != nil {
						logger.Error("Failed to create gzip reader", zap.Error(err))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					defer func(gr *gzipReader) {
						err := gr.Close()
						if err != nil {
							logger.Error("Failed to close gzip reader", zap.Error(err))
						}
					}(gr)
					r.Body = gr
				}

				next.ServeHTTP(ow, r)
			},
		)
	}
}
