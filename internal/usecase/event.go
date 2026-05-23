package usecase

import (
	"strings"

	"github.com/StepOne-ai/rwb_popular_requests/internal/metrics"
	"github.com/StepOne-ai/rwb_popular_requests/internal/schema"
)

func (u *Usecase) ProcessEvent(event *schema.SearchEvent) {
	query := normalize(event.Query)
	if query == "" || event.UserID == "" {
		metrics.EventsDropped.Inc()
		return
	}
	u.store.Increment(query, event.UserID, event.Timestamp)
	metrics.EventsProcessed.Inc()
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
