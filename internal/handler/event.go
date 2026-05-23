package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/StepOne-ai/rwb_popular_requests/internal/schema"
)

func (h *Handler) IngestEvent(w http.ResponseWriter, r *http.Request) {
	var event schema.SearchEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	if err := h.prod.Publish(r.Context(), &event); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to publish event")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
