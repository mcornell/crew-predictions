package repository

import (
	"context"
	"sync"
)

type Prediction struct {
	MatchID   string
	UserID    string
	Handle    string
	HomeGoals int
	AwayGoals int
}

type PredictionStore interface {
	Save(ctx context.Context, p Prediction) error
	GetByMatchAndUser(ctx context.Context, matchID, userID string) (*Prediction, error)
	GetByMatch(ctx context.Context, matchID string) ([]Prediction, error)
	GetAll(ctx context.Context) ([]Prediction, error)
}

type MemoryPredictionStore struct {
	mu   sync.RWMutex
	data map[string]Prediction
}

func NewMemoryPredictionStore() *MemoryPredictionStore {
	return &MemoryPredictionStore{data: make(map[string]Prediction)}
}

func (s *MemoryPredictionStore) Save(ctx context.Context, p Prediction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[p.MatchID+"|"+p.UserID] = p
	return nil
}

func (s *MemoryPredictionStore) GetByMatchAndUser(ctx context.Context, matchID, userID string) (*Prediction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.data[matchID+"|"+userID]
	if !ok {
		return nil, nil
	}
	return &p, nil
}

func (s *MemoryPredictionStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]Prediction)
}

func (s *MemoryPredictionStore) GetByMatch(ctx context.Context, matchID string) ([]Prediction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []Prediction
	for _, p := range s.data {
		if p.MatchID == matchID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (s *MemoryPredictionStore) GetAll(ctx context.Context) ([]Prediction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	all := make([]Prediction, 0, len(s.data))
	for _, p := range s.data {
		all = append(all, p)
	}
	return all, nil
}
