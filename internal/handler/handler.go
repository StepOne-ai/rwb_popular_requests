package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/StepOne-ai/rwb_popular_requests/internal/schema"
)

type IUsecase interface {
	GetTop(n int) schema.TopResponse
	AddStop(word string) bool
	RemoveStop(word string) bool
	ListStop() []string
}

type IProducer interface {
	Publish(ctx context.Context, event *schema.SearchEvent) error
}

type Handler struct {
	uc      IUsecase
	prod    IProducer
	maxTopN int
}

func New(uc IUsecase, prod IProducer, maxTopN int) *Handler {
	return &Handler{uc: uc, prod: prod, maxTopN: maxTopN}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, schema.ErrorResponse{Error: msg})
}
