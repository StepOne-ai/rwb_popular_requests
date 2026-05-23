package repository

import (
	"testing"
	"time"
)

func newTestStore() *Store {
	return NewStore(5, time.Minute)
}

func TestGetCachedTopEmpty(t *testing.T) {
	s := newTestStore()
	entries := s.GetCachedTop(10)
	if len(entries) != 0 {
		t.Fatalf("expected empty top on fresh store, got %d entries", len(entries))
	}
}

func TestRecomputeBasic(t *testing.T) {
	s := newTestStore()
	now := time.Now()

	s.Increment("кроссовки", "u1", now)
	s.Increment("кроссовки", "u2", now)
	s.Increment("платье", "u1", now)

	s.recompute()

	top := s.GetCachedTop(10)
	if len(top) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(top))
	}
	if top[0].Query != "кроссовки" || top[0].Count != 2 {
		t.Fatalf("expected кроссовки:2 at top, got %s:%d", top[0].Query, top[0].Count)
	}
	if top[1].Query != "платье" || top[1].Count != 1 {
		t.Fatalf("expected платье:1 at second, got %s:%d", top[1].Query, top[1].Count)
	}
}

func TestGetCachedTopLimitN(t *testing.T) {
	s := newTestStore()
	now := time.Now()

	queries := []string{"a", "b", "c", "d", "e"}
	for i, q := range queries {
		for j := 0; j <= i; j++ {
			s.Increment(q, "u"+string(rune('0'+j)), now)
		}
	}
	s.recompute()

	top := s.GetCachedTop(3)
	if len(top) != 3 {
		t.Fatalf("expected 3 entries with n=3, got %d", len(top))
	}
}

func TestStopListFiltersFromTop(t *testing.T) {
	s := newTestStore()
	now := time.Now()

	s.Increment("казино", "u1", now)
	s.Increment("казино", "u2", now)
	s.Increment("кроссовки", "u1", now)

	s.AddStop("казино")
	s.recompute()

	top := s.GetCachedTop(10)
	for _, e := range top {
		if e.Query == "казино" {
			t.Fatal("stop-list word appeared in top")
		}
	}
	if len(top) != 1 || top[0].Query != "кроссовки" {
		t.Fatalf("expected only кроссовки, got %+v", top)
	}
}

func TestAddRemoveStop(t *testing.T) {
	s := newTestStore()

	if !s.AddStop("казино") {
		t.Fatal("first add should return true")
	}
	if s.AddStop("казино") {
		t.Fatal("duplicate add should return false")
	}
	if !s.RemoveStop("казино") {
		t.Fatal("remove existing should return true")
	}
	if s.RemoveStop("казино") {
		t.Fatal("remove missing should return false")
	}
}

func TestRotateBucketClearsCounts(t *testing.T) {
	s := newTestStore()
	now := time.Now()

	s.Increment("кроссовки", "u1", now)
	s.recompute()

	// прокручиваем все 5 бакетов — старые данные должны уйти
	for i := 0; i < s.size; i++ {
		s.rotateBucket()
	}
	s.recompute()

	top := s.GetCachedTop(10)
	if len(top) != 0 {
		t.Fatalf("expected empty top after full rotation, got %d entries", len(top))
	}
}
