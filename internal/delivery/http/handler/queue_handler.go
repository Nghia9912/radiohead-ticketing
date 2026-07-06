package handler

import (
	"encoding/json"
	"net/http"

	"github.com/NghiaHoang/radiohead-ticketing/internal/repository/redis"
)

type QueueHandler struct {
	repo *redis.QueueRedisRepo
}

func NewQueueHandler(repo *redis.QueueRedisRepo) *QueueHandler {
	return &QueueHandler{repo: repo}
}

func (h *QueueHandler) Enqueue(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from query param (or use remote address as fallback for testing)
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "anonymous-" + r.RemoteAddr
	}

	err := h.repo.Enqueue(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to enqueue", http.StatusInternalServerError)
		return
	}

	rank, err := h.repo.GetRank(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get rank", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Enqueued successfully",
		"user_id": userID,
		"rank":    rank,
	})
}
