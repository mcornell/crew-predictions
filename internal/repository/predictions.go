package repository

import "context"

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
	GetAll(ctx context.Context) ([]Prediction, error)
}

type MemoryPredictionStore struct {
	data map[string]Prediction
}

func NewMemoryPredictionStore() *MemoryPredictionStore {
	return &MemoryPredictionStore{data: make(map[string]Prediction)}
}

func (s *MemoryPredictionStore) Save(ctx context.Context, p Prediction) error {
	s.data[p.MatchID+"|"+p.UserID] = p
	return nil
}

func (s *MemoryPredictionStore) GetByMatchAndUser(ctx context.Context, matchID, userID string) (*Prediction, error) {
	p, ok := s.data[matchID+"|"+userID]
	if !ok {
		return nil, nil
	}
	return &p, nil
}

func (s *MemoryPredictionStore) GetAll(ctx context.Context) ([]Prediction, error) {
	all := make([]Prediction, 0, len(s.data))
	for _, p := range s.data {
		all = append(all, p)
	}
	return all, nil
}
