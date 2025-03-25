package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"test_task/internal/service"
	"time"
)

type ReqBody struct {
	TsFrom string `json:"TsFrom"`
	TsTo   string `json:"TsTo"`
}

func NewStatsHandler(service *service.ClickStats) http.HandlerFunc {
	const timeLayout = "2006-01-02 15:04:05"
	return func(w http.ResponseWriter, r *http.Request) {
		bannerIDStr := r.URL.Path[len("/stats/"):]
		bannerID, err := strconv.Atoi(bannerIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var req ReqBody
		err = json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		tsFrom, err := time.Parse(timeLayout, req.TsFrom)
		if err != nil {
			http.Error(w, "invalid from timestamp", http.StatusBadRequest)
			return
		}
		tsTo, err := time.Parse(timeLayout, req.TsTo)
		if err != nil {
			http.Error(w, "invalid to timestamp", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		result, err := service.GetClickStats(ctx, bannerID, tsFrom, tsTo)
		if err != nil {
			slog.Info("something went wrong", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result:"` + string(result) + `}`))
	}
}
