package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/StepOne-ai/rwb_popular_requests/internal/config"
	"github.com/StepOne-ai/rwb_popular_requests/internal/consumer"
	"github.com/StepOne-ai/rwb_popular_requests/internal/handler"
	"github.com/StepOne-ai/rwb_popular_requests/internal/metrics"
	"github.com/StepOne-ai/rwb_popular_requests/internal/producer"
	"github.com/StepOne-ai/rwb_popular_requests/internal/repository"
	"github.com/StepOne-ai/rwb_popular_requests/internal/usecase"
)

func main() {
	cfg := config.Load()

	store := repository.NewStore(cfg.WindowSize, cfg.BucketDuration)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	store.Run(ctx, cfg.BucketDuration, cfg.CacheRefresh)

	uc := usecase.New(store)
	prod := producer.New(cfg.KafkaBrokers, cfg.KafkaTopic)
	defer prod.Close()
	h := handler.New(uc, prod, cfg.MaxTopN)

	cons := consumer.New(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroup, uc)
	go func() {
		if err := cons.Run(ctx); err != nil {
			log.Printf("consumer stopped: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/event", h.IngestEvent)
	mux.HandleFunc("GET /api/v1/top", h.GetTop)
	mux.HandleFunc("GET /api/v1/stoplist", h.ListStop)
	mux.HandleFunc("POST /api/v1/stoplist", h.AddStop)
	mux.HandleFunc("DELETE /api/v1/stoplist/{word}", h.RemoveStop)
	mux.Handle("GET /metrics", metrics.Handler())
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("http server listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	shutCtx, shutCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutCancel()
	_ = srv.Shutdown(shutCtx)
}
