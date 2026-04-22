package repository

import (
	"context"
	"sync"
)

type User struct {
	UserID        string
	Handle        string
	Provider      string
	EmailVerified bool
}

type UserStore interface {
	Upsert(ctx context.Context, u User) error
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
	s.data[u.UserID] = u
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
