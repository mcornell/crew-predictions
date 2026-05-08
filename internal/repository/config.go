// Package repository defines the storage interfaces used across the server
// and provides both in-memory and Firestore-backed implementations for
// matches, predictions, results, users, seasons, and global config state.
package repository

import (
	"context"
	"sync"
)

type ConfigStore interface {
	GetActiveSeason(ctx context.Context) string
	SetActiveSeason(ctx context.Context, id string) error
}

type MemoryConfigStore struct {
	mu        sync.RWMutex
	seasonID  string
	defaultID string
}

func NewMemoryConfigStore(defaultSeasonID string) *MemoryConfigStore {
	return &MemoryConfigStore{seasonID: defaultSeasonID, defaultID: defaultSeasonID}
}

func (c *MemoryConfigStore) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.seasonID = c.defaultID
}

func (c *MemoryConfigStore) GetActiveSeason(_ context.Context) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.seasonID
}

func (c *MemoryConfigStore) SetActiveSeason(_ context.Context, id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.seasonID = id
	return nil
}
