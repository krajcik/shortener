package handler

import (
	"context"
	"database/sql"
	"net/http"

	"go.uber.org/zap"
)

func Ping(db *sql.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//duration, err := time.ParseDuration("1s")
		//if err != nil {
		//	logger.Error("Failed to parse duration", zap.Error(err))
		//}
		//ctx, cancelFunc := context.WithTimeout(context.Background(), duration)
		//defer cancelFunc()
		err := db.PingContext(context.Background())
		if err != nil {
			logger.Error("Failed to ping database", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
