package repository

import (
	"context"
	"sync"
	"time"
)

type SeasonEntry struct {
	UserID          string `firestore:"userID"          json:"userID"`
	Handle          string `firestore:"handle"          json:"handle"`
	AcesRadioPoints int    `firestore:"acesRadioPoints" json:"acesRadioPoints"`
	Upper90Points   int    `firestore:"upper90Points"   json:"upper90Points"`
	GrouchyPoints   int    `firestore:"grouchyPoints"   json:"grouchyPoints"`
	PredictionCount int    `firestore:"predictionCount" json:"predictionCount"`
	Rank            int    `firestore:"rank"            json:"rank"`
}

type SeasonSnapshot struct {
	ID       string        `firestore:"id"`
	Name     string        `firestore:"name"`
	ClosedAt time.Time     `firestore:"closedAt"`
	Entries  []SeasonEntry `firestore:"entries"`
}

type SeasonStore interface {
	Save(ctx context.Context, s SeasonSnapshot) error
	GetByID(ctx context.Context, id string) (*SeasonSnapshot, error)
	ListAll(ctx context.Context) ([]SeasonSnapshot, error)
	Exists(ctx context.Context, id string) bool
	Reset()
}

type MemorySeasonStore struct {
	mu   sync.RWMutex
	data map[string]SeasonSnapshot
}

func NewMemorySeasonStore() *MemorySeasonStore {
	return &MemorySeasonStore{data: make(map[string]SeasonSnapshot)}
}

func (s *MemorySeasonStore) Save(_ context.Context, snap SeasonSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[snap.ID] = snap
	return nil
}

func (s *MemorySeasonStore) GetByID(_ context.Context, id string) (*SeasonSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap, ok := s.data[id]
	if !ok {
		return nil, nil
	}
	return &snap, nil
}

func (s *MemorySeasonStore) ListAll(_ context.Context) ([]SeasonSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]SeasonSnapshot, 0, len(s.data))
	for _, snap := range s.data {
		out = append(out, snap)
	}
	return out, nil
}

func (s *MemorySeasonStore) Exists(_ context.Context, id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[id]
	return ok
}

func (s *MemorySeasonStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]SeasonSnapshot)
}
