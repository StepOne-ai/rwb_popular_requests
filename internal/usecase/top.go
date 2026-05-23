package usecase

import (
	"time"

	"github.com/StepOne-ai/rwb_popular_requests/internal/schema"
)

func (u *Usecase) GetTop(n int) schema.TopResponse {
	entries := u.store.GetCachedTop(n)
	return schema.TopResponse{
		Items:         entries,
		WindowMinutes: 5,
		GeneratedAt:   time.Now().UTC(),
	}
}
