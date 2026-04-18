package repository

import "context"

type Result struct {
	MatchID   string
	HomeGoals int
	AwayGoals int
}

type ResultStore interface {
	SaveResult(ctx context.Context, r Result) error
	GetResult(ctx context.Context, matchID string) (*Result, error)
}

type MemoryResultStore struct {
	data map[string]Result
}

func NewMemoryResultStore() *MemoryResultStore {
	return &MemoryResultStore{data: make(map[string]Result)}
}

func (s *MemoryResultStore) SaveResult(ctx context.Context, r Result) error {
	s.data[r.MatchID] = r
	return nil
}

func (s *MemoryResultStore) GetResult(ctx context.Context, matchID string) (*Result, error) {
	r, ok := s.data[matchID]
	if !ok {
		return nil, nil
	}
	return &r, nil
}
