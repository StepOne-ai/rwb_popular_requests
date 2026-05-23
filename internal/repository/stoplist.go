package repository

import "sync"

type stopList struct {
	mu    sync.RWMutex
	words map[string]struct{}
}

func newStopList() *stopList {
	return &stopList{words: make(map[string]struct{})}
}

func (s *stopList) add(word string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.words[word]; exists {
		return false
	}
	s.words[word] = struct{}{}
	return true
}

func (s *stopList) remove(word string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.words[word]; !exists {
		return false
	}
	delete(s.words, word)
	return true
}

func (s *stopList) contains(word string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.words[word]
	return ok
}

func (s *stopList) list() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]string, 0, len(s.words))
	for w := range s.words {
		out = append(out, w)
	}
	return out
}
