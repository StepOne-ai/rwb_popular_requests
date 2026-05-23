package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	EventsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "popular_requests_events_processed_total",
		Help: "Total number of search events processed",
	})

	EventsDropped = promauto.NewCounter(prometheus.CounterOpts{
		Name: "popular_requests_events_dropped_total",
		Help: "Total number of search events dropped (empty query or user_id)",
	})

	TopRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "popular_requests_top_requests_total",
		Help: "Total number of GET /top requests",
	})

	StopListSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "popular_requests_stoplist_size",
		Help: "Current number of words in the stop-list",
	})
)

func Handler() http.Handler {
	return promhttp.Handler()
}
