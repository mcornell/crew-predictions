package repository

import (
	"context"
	"sync"
)

type Result struct {
	MatchID   string
	HomeTeam  string
	AwayTeam  string
	HomeGoals int
	AwayGoals int
}

type ResultStore interface {
	SaveResult(ctx context.Context, r Result) error
	GetResult(ctx context.Context, matchID string) (*Result, error)
}

type MemoryResultStore struct {
	mu   sync.RWMutex
	data map[string]Result
}

func NewMemoryResultStore() *MemoryResultStore {
	return &MemoryResultStore{data: make(map[string]Result)}
}

func (s *MemoryResultStore) SaveResult(ctx context.Context, r Result) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[r.MatchID] = r
	return nil
}

func (s *MemoryResultStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]Result)
}

func (s *MemoryResultStore) GetResult(ctx context.Context, matchID string) (*Result, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.data[matchID]
	if !ok {
		return nil, nil
	}
	return &r, nil
}
