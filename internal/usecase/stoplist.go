package usecase

import (
	"strings"

	"github.com/StepOne-ai/rwb_popular_requests/internal/metrics"
)

func (u *Usecase) AddStop(word string) bool {
	ok := u.store.AddStop(strings.ToLower(strings.TrimSpace(word)))
	if ok {
		metrics.StopListSize.Inc()
	}
	return ok
}

func (u *Usecase) RemoveStop(word string) bool {
	ok := u.store.RemoveStop(strings.ToLower(strings.TrimSpace(word)))
	if ok {
		metrics.StopListSize.Dec()
	}
	return ok
}

func (u *Usecase) ListStop() []string {
	return u.store.ListStop()
}
