package repository

import "context"

type Prediction struct {
	MatchID   string
	Handle    string
	HomeGoals int
	AwayGoals int
}

type PredictionStore interface {
	Save(ctx context.Context, p Prediction) error
	GetByMatchAndHandle(ctx context.Context, matchID, handle string) (*Prediction, error)
	GetAll(ctx context.Context) ([]Prediction, error)
}

type MemoryPredictionStore struct {
	data map[string]Prediction
}

func NewMemoryPredictionStore() *MemoryPredictionStore {
	return &MemoryPredictionStore{data: make(map[string]Prediction)}
}

func (s *MemoryPredictionStore) Save(ctx context.Context, p Prediction) error {
	s.data[p.MatchID+"|"+p.Handle] = p
	return nil
}

func (s *MemoryPredictionStore) GetByMatchAndHandle(ctx context.Context, matchID, handle string) (*Prediction, error) {
	p, ok := s.data[matchID+"|"+handle]
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
