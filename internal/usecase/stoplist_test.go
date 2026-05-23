package usecase

import (
	"testing"
)

func TestAddStop_NormalizesWord(t *testing.T) {
	store := newMockStore()
	uc := New(store)

	uc.AddStop("  КАЗИНО  ")

	if _, ok := store.stopList["казино"]; !ok {
		t.Fatal("expected 'казино' in stoplist after adding '  КАЗИНО  '")
	}
}

func TestAddStop_Duplicate(t *testing.T) {
	store := newMockStore()
	uc := New(store)

	if !uc.AddStop("казино") {
		t.Fatal("first add should return true")
	}
	if uc.AddStop("казино") {
		t.Fatal("duplicate add should return false")
	}
}

func TestRemoveStop_Existing(t *testing.T) {
	store := newMockStore()
	uc := New(store)

	uc.AddStop("казино")
	if !uc.RemoveStop("казино") {
		t.Fatal("remove existing word should return true")
	}
	if _, ok := store.stopList["казино"]; ok {
		t.Fatal("word should be gone after remove")
	}
}

func TestRemoveStop_NotFound(t *testing.T) {
	store := newMockStore()
	uc := New(store)

	if uc.RemoveStop("несуществующее") {
		t.Fatal("remove missing word should return false")
	}
}

func TestListStop(t *testing.T) {
	store := newMockStore()
	uc := New(store)

	uc.AddStop("казино")
	uc.AddStop("ставки")

	words := uc.ListStop()
	if len(words) != 2 {
		t.Fatalf("expected 2 words, got %d", len(words))
	}
}
