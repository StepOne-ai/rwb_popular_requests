package usecase

import (
	"sync"
	"time"

	"github.com/StepOne-ai/rwb_popular_requests/internal/schema"
)

type mockStore struct {
	mu       sync.Mutex
	calls    []incrementCall
	top      []schema.TopEntry
	stopList map[string]struct{}
}

type incrementCall struct {
	query  string
	userID string
}

func newMockStore() *mockStore {
	return &mockStore{stopList: make(map[string]struct{})}
}

func (m *mockStore) Increment(query, userID string, _ time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, incrementCall{query, userID})
}

func (m *mockStore) GetCachedTop(n int) []schema.TopEntry {
	if n >= len(m.top) {
		return m.top
	}
	return m.top[:n]
}

func (m *mockStore) AddStop(word string) bool {
	if _, ok := m.stopList[word]; ok {
		return false
	}
	m.stopList[word] = struct{}{}
	return true
}

func (m *mockStore) RemoveStop(word string) bool {
	if _, ok := m.stopList[word]; !ok {
		return false
	}
	delete(m.stopList, word)
	return true
}

func (m *mockStore) ListStop() []string {
	out := make([]string, 0, len(m.stopList))
	for w := range m.stopList {
		out = append(out, w)
	}
	return out
}
