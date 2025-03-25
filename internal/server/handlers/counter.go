package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"test_task/internal/service"
)

func NewClickHandler(cc *service.ClickCounter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bannerIDStr := strings.TrimPrefix(r.URL.Path, "/counter/")
		if bannerIDStr == "" {
			http.Error(w, "Empty banner ID", http.StatusBadRequest)
			return
		}

		bannerID, err := strconv.Atoi(bannerIDStr)
		if err != nil {
			http.Error(w, "Invalid banner ID format", http.StatusBadRequest)
			return
		}

		if bannerID < 1 {
			http.Error(w, "Banner ID must be positive", http.StatusBadRequest)
			return
		}

		cc.AddClick(bannerID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status: OK"}`))
	}
}
