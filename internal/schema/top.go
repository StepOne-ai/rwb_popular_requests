package schema

import "time"

type TopEntry struct {
	Query string `json:"query"`
	Count uint64 `json:"count"`
}

type TopResponse struct {
	Items         []TopEntry `json:"items"`
	WindowMinutes int        `json:"window_minutes"`
	GeneratedAt   time.Time  `json:"generated_at"`
}
