package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr string

	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroup   string

	WindowSize      int           // кол-во бакетов (= минут)
	BucketDuration  time.Duration
	CacheRefresh    time.Duration // как часто рефрешить top-N

	DefaultTopN int
	MaxTopN     int
}

func Load() *Config {
	return &Config{
		HTTPAddr:       getEnv("HTTP_ADDR", ":8080"),
		KafkaBrokers:   []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
		KafkaTopic:     getEnv("KAFKA_TOPIC", "search-events"),
		KafkaGroup:     getEnv("KAFKA_GROUP", "popular-requests"),
		WindowSize:     5,
		BucketDuration: time.Minute,
		CacheRefresh:   5 * time.Second,
		DefaultTopN:    10,
		MaxTopN:        getEnvInt("MAX_TOP_N", 100),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
