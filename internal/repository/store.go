package repository

import (
	"context"
	"sort"
	"sync/atomic"
	"time"

	"github.com/StepOne-ai/rwb_popular_requests/internal/schema"
)

type Store struct {
	buckets []*bucket
	size    int

	current atomic.Int32
	cached  atomic.Pointer[[]schema.TopEntry]

	stopList *stopList
}

func NewStore(windowSize int, bucketDuration time.Duration) *Store {
	buckets := make([]*bucket, windowSize)
	now := time.Now()
	for i := range buckets {
		buckets[i] = newBucket(now.Add(-time.Duration(windowSize-i) * bucketDuration))
	}

	s := &Store{
		buckets:  buckets,
		size:     windowSize,
		stopList: newStopList(),
	}

	empty := []schema.TopEntry{}
	s.cached.Store(&empty)
	return s
}

func (s *Store) Increment(query, userID string, t time.Time) {
	idx := s.current.Load()
	cur := s.buckets[idx]

	// если событие попало в предыдущую минуту — пишем в предыдущий бакет
	if t.Before(cur.start) {
		prev := (int(idx) - 1 + s.size) % s.size
		if t.After(s.buckets[prev].start) {
			s.buckets[prev].increment(query, userID)
			return
		}
	}
	cur.increment(query, userID)
}

// recompute сливает все бакеты, применяет стоп-лист, сортирует и кеширует.
func (s *Store) recompute() {
	merged := make(map[string]uint64)
	for _, b := range s.buckets {
		for q, c := range b.snapshot() {
			if !s.stopList.contains(q) {
				merged[q] += c
			}
		}
	}

	entries := make([]schema.TopEntry, 0, len(merged))
	for q, c := range merged {
		entries = append(entries, schema.TopEntry{Query: q, Count: c})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Count > entries[j].Count
	})

	s.cached.Store(&entries)
}

func (s *Store) GetCachedTop(n int) []schema.TopEntry {
	all := *s.cached.Load()
	if n >= len(all) {
		return all
	}
	return all[:n]
}

// rotateBucket сбрасывает самый старый бакет и делает его новым текущим.
func (s *Store) rotateBucket() {
	next := (int(s.current.Load()) + 1) % s.size
	s.buckets[next].reset(time.Now())
	s.current.Store(int32(next))
}

func (s *Store) AddStop(word string) bool    { return s.stopList.add(word) }
func (s *Store) RemoveStop(word string) bool { return s.stopList.remove(word) }
func (s *Store) ListStop() []string          { return s.stopList.list() }

// Run запускает фоновые горутины ротации и пересчёта.
func (s *Store) Run(ctx context.Context, bucketDuration, cacheRefresh time.Duration) {
	go func() {
		t := time.NewTicker(bucketDuration)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				s.rotateBucket()
			}
		}
	}()

	go func() {
		t := time.NewTicker(cacheRefresh)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				s.recompute()
			}
		}
	}()
}
