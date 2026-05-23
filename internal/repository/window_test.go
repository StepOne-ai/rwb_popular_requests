package repository

import (
	"testing"
	"time"
)

func TestBucketDedup(t *testing.T) {
	b := newBucket(time.Now())

	b.increment("кроссовки", "user1")
	b.increment("кроссовки", "user1")
	b.increment("кроссовки", "user1")

	if got := b.counts["кроссовки"]; got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestBucketDifferentUsers(t *testing.T) {
	b := newBucket(time.Now())

	b.increment("платье", "user1")
	b.increment("платье", "user2")
	b.increment("платье", "user3")

	if got := b.counts["платье"]; got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestBucketSameUserDifferentQueries(t *testing.T) {
	b := newBucket(time.Now())

	b.increment("кроссовки", "user1")
	b.increment("платье", "user1")
	b.increment("наушники", "user1")

	for _, q := range []string{"кроссовки", "платье", "наушники"} {
		if got := b.counts[q]; got != 1 {
			t.Fatalf("query %q: expected 1, got %d", q, got)
		}
	}
}

func TestBucketReset(t *testing.T) {
	b := newBucket(time.Now())
	b.increment("кроссовки", "user1")
	b.increment("платье", "user2")

	b.reset(time.Now())

	if len(b.counts) != 0 {
		t.Fatalf("expected empty counts after reset, got %d entries", len(b.counts))
	}
	if len(b.seen) != 0 {
		t.Fatalf("expected empty seen after reset, got %d entries", len(b.seen))
	}
}

func TestBucketSnapshot(t *testing.T) {
	b := newBucket(time.Now())
	b.increment("кроссовки", "user1")
	b.increment("кроссовки", "user2")

	snap := b.snapshot()

	// мутируем оригинал — снимок не должен измениться
	b.increment("кроссовки", "user3")

	if got := snap["кроссовки"]; got != 2 {
		t.Fatalf("snapshot: expected 2, got %d", got)
	}
}
