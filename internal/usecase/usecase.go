package usecase

import (
	"time"

	"github.com/StepOne-ai/rwb_popular_requests/internal/schema"
)

type IStore interface {
	Increment(query, userID string, t time.Time)
	GetCachedTop(n int) []schema.TopEntry
	AddStop(word string) bool
	RemoveStop(word string) bool
	ListStop() []string
}

type Usecase struct {
	store IStore
}

func New(store IStore) *Usecase {
	return &Usecase{store: store}
}
