package handler

import (
	"net/http"
	"strconv"

	"github.com/StepOne-ai/rwb_popular_requests/internal/metrics"
)

func (h *Handler) GetTop(w http.ResponseWriter, r *http.Request) {
	n := 10
	if raw := r.URL.Query().Get("n"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 1 || v > h.maxTopN {
			writeError(w, http.StatusBadRequest, "n must be an integer between 1 and "+strconv.Itoa(h.maxTopN))
			return
		}
		n = v
	}

	metrics.TopRequests.Inc()
	resp := h.uc.GetTop(n)
	writeJSON(w, http.StatusOK, resp)
}
