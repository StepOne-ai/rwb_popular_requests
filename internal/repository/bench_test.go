package repository

import (
	"fmt"
	"testing"
	"time"
)

// BenchmarkIncrement — скорость записи уникальных событий в бакет.
func BenchmarkIncrement(b *testing.B) {
	s := NewStore(5, time.Minute)
	now := time.Now()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			s.Increment("кроссовки", fmt.Sprintf("user-%d", i%10000), now)
			i++
		}
	})
}

// BenchmarkIncrementUniqueQueries — много разных запросов, реалистичная нагрузка.
func BenchmarkIncrementUniqueQueries(b *testing.B) {
	s := NewStore(5, time.Minute)
	now := time.Now()
	queries := make([]string, 1000)
	for i := range queries {
		queries[i] = fmt.Sprintf("query-%d", i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			s.Increment(queries[i%len(queries)], fmt.Sprintf("user-%d", i%10000), now)
			i++
		}
	})
}

// BenchmarkGetCachedTop — скорость чтения топа, это горячий путь.
func BenchmarkGetCachedTop(b *testing.B) {
	s := NewStore(5, time.Minute)
	now := time.Now()

	// наполняем данными
	for i := range 1000 {
		s.Increment(fmt.Sprintf("query-%d", i), fmt.Sprintf("user-%d", i), now)
	}
	s.recompute()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = s.GetCachedTop(10)
		}
	})
}

// BenchmarkRecompute — скорость фонового пересчёта топа.
func BenchmarkRecompute(b *testing.B) {
	s := NewStore(5, time.Minute)
	now := time.Now()

	for i := range 10000 {
		s.Increment(fmt.Sprintf("query-%d", i%500), fmt.Sprintf("user-%d", i), now)
	}

	b.ResetTimer()
	for b.Loop() {
		s.recompute()
	}
}
