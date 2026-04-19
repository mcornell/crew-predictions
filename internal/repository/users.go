package repository

import "context"

type User struct {
	UserID   string
	Handle   string
	Provider string
}

type UserStore interface {
	Upsert(ctx context.Context, u User) error
	GetUser(ctx context.Context, userID string) (*User, error)
}

type MemoryUserStore struct {
	data map[string]User
}

func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{data: make(map[string]User)}
}

func (s *MemoryUserStore) Upsert(ctx context.Context, u User) error {
	s.data[u.UserID] = u
	return nil
}

func (s *MemoryUserStore) GetUser(ctx context.Context, userID string) (*User, error) {
	u, ok := s.data[userID]
	if !ok {
		return nil, nil
	}
	return &u, nil
}
