package usecase

import (
	"testing"
	"time"

	"github.com/StepOne-ai/rwb_popular_requests/internal/schema"
)

func TestProcessEvent_NormalizesQuery(t *testing.T) {
	store := newMockStore()
	uc := New(store)

	uc.ProcessEvent(&schema.SearchEvent{
		Query:     "  КРОССОВКИ  ",
		UserID:    "u1",
		Timestamp: time.Now(),
	})

	if len(store.calls) != 1 {
		t.Fatalf("expected 1 increment call, got %d", len(store.calls))
	}
	if store.calls[0].query != "кроссовки" {
		t.Fatalf("expected normalized query 'кроссовки', got %q", store.calls[0].query)
	}
}

func TestProcessEvent_DropsEmptyQuery(t *testing.T) {
	store := newMockStore()
	uc := New(store)

	uc.ProcessEvent(&schema.SearchEvent{
		Query:     "   ",
		UserID:    "u1",
		Timestamp: time.Now(),
	})

	if len(store.calls) != 0 {
		t.Fatalf("expected 0 increment calls for empty query, got %d", len(store.calls))
	}
}

func TestProcessEvent_DropsEmptyUserID(t *testing.T) {
	store := newMockStore()
	uc := New(store)

	uc.ProcessEvent(&schema.SearchEvent{
		Query:     "кроссовки",
		UserID:    "",
		Timestamp: time.Now(),
	})

	if len(store.calls) != 0 {
		t.Fatalf("expected 0 increment calls for empty user_id, got %d", len(store.calls))
	}
}

func TestProcessEvent_PassesCorrectFields(t *testing.T) {
	store := newMockStore()
	uc := New(store)

	uc.ProcessEvent(&schema.SearchEvent{
		Query:     "наушники",
		UserID:    "u42",
		Timestamp: time.Now(),
	})

	if store.calls[0].userID != "u42" {
		t.Fatalf("expected userID 'u42', got %q", store.calls[0].userID)
	}
}
