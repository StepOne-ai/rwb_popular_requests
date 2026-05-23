package usecase

import (
	"testing"
	"time"

	"github.com/StepOne-ai/rwb_popular_requests/internal/schema"
)

func TestGetTop_ReturnsStoreEntries(t *testing.T) {
	store := newMockStore()
	store.top = []schema.TopEntry{
		{Query: "кроссовки", Count: 100},
		{Query: "платье", Count: 50},
	}
	uc := New(store)

	resp := uc.GetTop(10)

	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}
	if resp.Items[0].Query != "кроссовки" || resp.Items[0].Count != 100 {
		t.Fatalf("unexpected first item: %+v", resp.Items[0])
	}
}

func TestGetTop_RespectsN(t *testing.T) {
	store := newMockStore()
	store.top = []schema.TopEntry{
		{Query: "a", Count: 3},
		{Query: "b", Count: 2},
		{Query: "c", Count: 1},
	}
	uc := New(store)

	resp := uc.GetTop(2)

	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items with n=2, got %d", len(resp.Items))
	}
}

func TestGetTop_EmptyStore(t *testing.T) {
	store := newMockStore()
	uc := New(store)

	resp := uc.GetTop(10)

	if len(resp.Items) != 0 {
		t.Fatalf("expected empty items, got %d", len(resp.Items))
	}
}

func TestGetTop_MetaFields(t *testing.T) {
	store := newMockStore()
	uc := New(store)

	before := time.Now().UTC()
	resp := uc.GetTop(10)
	after := time.Now().UTC()

	if resp.WindowMinutes != 5 {
		t.Fatalf("expected window_minutes=5, got %d", resp.WindowMinutes)
	}
	if resp.GeneratedAt.Before(before) || resp.GeneratedAt.After(after) {
		t.Fatalf("generated_at out of range: %v", resp.GeneratedAt)
	}
}
