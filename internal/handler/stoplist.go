package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

type stopWordRequest struct {
	Word string `json:"word"`
}

func (h *Handler) AddStop(w http.ResponseWriter, r *http.Request) {
	var req stopWordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Word) == "" {
		writeError(w, http.StatusBadRequest, "word is required")
		return
	}

	if !h.uc.AddStop(req.Word) {
		writeError(w, http.StatusConflict, "already in stoplist")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"word": req.Word, "action": "added"})
}

func (h *Handler) RemoveStop(w http.ResponseWriter, r *http.Request) {
	word := r.PathValue("word")
	if word == "" {
		writeError(w, http.StatusBadRequest, "word is required")
		return
	}

	if !h.uc.RemoveStop(word) {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"word": word, "action": "removed"})
}

func (h *Handler) ListStop(w http.ResponseWriter, r *http.Request) {
	words := h.uc.ListStop()
	writeJSON(w, http.StatusOK, map[string]any{"words": words, "total": len(words)})
}
