package repository

import (
	"maps"
	"sync"
	"time"
)

type bucket struct {
	mu     sync.Mutex
	counts map[string]uint64          // query → уник юзеров
	seen   map[string]map[string]bool // query → set of user_ids
	start  time.Time
}

func newBucket(start time.Time) *bucket {
	return &bucket{
		counts: make(map[string]uint64),
		seen:   make(map[string]map[string]bool),
		start:  start,
	}
}

func (b *bucket) increment(query, userID string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.seen[query]; !ok {
		b.seen[query] = make(map[string]bool)
	}
	if b.seen[query][userID] {
		return
	}
	b.seen[query][userID] = true
	b.counts[query]++
}

// snapshot возвращает копию counts без блокировки снаружи.
func (b *bucket) snapshot() map[string]uint64 {
	b.mu.Lock()
	defer b.mu.Unlock()

	out := make(map[string]uint64, len(b.counts))
	maps.Copy(out, b.counts)
	return out
}

func (b *bucket) reset(start time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.counts = make(map[string]uint64)
	b.seen = make(map[string]map[string]bool)
	b.start = start
}
