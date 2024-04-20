package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Logger(logger *zap.Logger, name string) func(next http.Handler) http.Handler {
	logger = zap.New(logger.Core(), zap.AddCallerSkip(1)).Named(name)
	logger.Debug("zap.logger detected for chi")
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				logger.Info("served",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Int("status", ww.Status()),
					zap.String("reqId", middleware.GetReqID(r.Context())),
					zap.String("remoteAddr", r.RemoteAddr),
					zap.String("proto", r.Proto),
					zap.Duration("latency", time.Since(t1)),
					zap.Int("size", ww.BytesWritten()))
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
