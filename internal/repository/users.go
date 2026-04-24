package repository

import (
	"context"
	"sync"
)

type User struct {
	UserID           string
	Handle           string
	Provider         string
	Location         string
	EmailVerified    bool
	AcesRadioPoints  int
	Upper90Points    int
	GrouchyPoints    int
	PredictionCount  int
}

type UserStore interface {
	Upsert(ctx context.Context, u User) error
	UpdateScores(ctx context.Context, userID string, count, aces, u90, grouchy int) error
	GetByID(ctx context.Context, userID string) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
}

type MemoryUserStore struct {
	mu   sync.RWMutex
	data map[string]User
}

func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{data: make(map[string]User)}
}

func (s *MemoryUserStore) Upsert(_ context.Context, u User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing := s.data[u.UserID]
	if u.Provider == "" {
		u.Provider = existing.Provider
	}
	if u.Location == "" {
		u.Location = existing.Location
	}
	u.AcesRadioPoints = existing.AcesRadioPoints
	u.Upper90Points = existing.Upper90Points
	u.GrouchyPoints = existing.GrouchyPoints
	u.PredictionCount = existing.PredictionCount
	s.data[u.UserID] = u
	return nil
}

func (s *MemoryUserStore) UpdateScores(_ context.Context, userID string, count, aces, u90, grouchy int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	u := s.data[userID]
	u.UserID = userID
	u.AcesRadioPoints = aces
	u.Upper90Points = u90
	u.GrouchyPoints = grouchy
	u.PredictionCount = count
	s.data[userID] = u
	return nil
}

func (s *MemoryUserStore) GetByID(_ context.Context, userID string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.data[userID]
	if !ok {
		return nil, nil
	}
	return &u, nil
}

func (s *MemoryUserStore) GetAll(_ context.Context) ([]User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	all := make([]User, 0, len(s.data))
	for _, u := range s.data {
		all = append(all, u)
	}
	return all, nil
}
