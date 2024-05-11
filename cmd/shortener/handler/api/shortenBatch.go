package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func (h *ShortenBatchHandler) ShortenBatch(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var req []ShortenBatchRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	var urls []string
	for _, reqItem := range req {
		urls = append(urls, reqItem.OriginalURL)
	}
	err = h.S.ShrtBatch(context.Background(), urls)
	if err != nil {
		panic(err)
	}
}

type ShortenBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
