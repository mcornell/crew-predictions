package repository

import (
	"context"
	"fmt"
)

// ErrorPredictionStore is a test double that always fails on Save.
type ErrorPredictionStore struct{}

func NewErrorPredictionStore() *ErrorPredictionStore { return &ErrorPredictionStore{} }

func (e *ErrorPredictionStore) Save(_ context.Context, _ Prediction) error {
	return fmt.Errorf("simulated store failure")
}

func (e *ErrorPredictionStore) GetByMatchAndUser(_ context.Context, _, _ string) (*Prediction, error) {
	return nil, nil
}

func (e *ErrorPredictionStore) GetByMatch(_ context.Context, _ string) ([]Prediction, error) {
	return nil, nil
}

func (e *ErrorPredictionStore) GetAll(_ context.Context) ([]Prediction, error) {
	return nil, nil
}

// GetByMatchAndUserErrorPredictionStore delegates to a MemoryPredictionStore for all
// methods except GetByMatchAndUser, which always fails. This lets tests verify that
// the bot skips prediction saves without losing access to GetByMatch/GetAll.
type GetByMatchAndUserErrorPredictionStore struct {
	*MemoryPredictionStore
}

func NewGetByMatchAndUserErrorPredictionStore() *GetByMatchAndUserErrorPredictionStore {
	return &GetByMatchAndUserErrorPredictionStore{NewMemoryPredictionStore()}
}

func (e *GetByMatchAndUserErrorPredictionStore) GetByMatchAndUser(_ context.Context, _, _ string) (*Prediction, error) {
	return nil, fmt.Errorf("simulated GetByMatchAndUser failure")
}

// ErrorGetAllPredictionStore is a test double that always fails on GetAll.
type ErrorGetAllPredictionStore struct{}

func NewErrorGetAllPredictionStore() *ErrorGetAllPredictionStore {
	return &ErrorGetAllPredictionStore{}
}

func (e *ErrorGetAllPredictionStore) Save(_ context.Context, _ Prediction) error { return nil }

func (e *ErrorGetAllPredictionStore) GetByMatchAndUser(_ context.Context, _, _ string) (*Prediction, error) {
	return nil, nil
}

func (e *ErrorGetAllPredictionStore) GetByMatch(_ context.Context, _ string) ([]Prediction, error) {
	return nil, nil
}

func (e *ErrorGetAllPredictionStore) GetAll(_ context.Context) ([]Prediction, error) {
	return nil, fmt.Errorf("simulated GetAll failure")
}

// ErrorUpsertUserStore is a test double that always fails on Upsert.
type ErrorUpsertUserStore struct{}

func NewErrorUpsertUserStore() *ErrorUpsertUserStore { return &ErrorUpsertUserStore{} }

func (e *ErrorUpsertUserStore) Upsert(_ context.Context, _ User) error {
	return fmt.Errorf("simulated upsert failure")
}

func (e *ErrorUpsertUserStore) GetByID(_ context.Context, _ string) (*User, error) {
	return nil, nil
}

func (e *ErrorUpsertUserStore) GetAll(_ context.Context) ([]User, error) {
	return nil, nil
}

// ErrorGetByIDUserStore is a test double that always fails on GetByID.
type ErrorGetByIDUserStore struct{}

func NewErrorGetByIDUserStore() *ErrorGetByIDUserStore { return &ErrorGetByIDUserStore{} }

func (e *ErrorGetByIDUserStore) Upsert(_ context.Context, _ User) error { return nil }

func (e *ErrorGetByIDUserStore) GetByID(_ context.Context, _ string) (*User, error) {
	return nil, fmt.Errorf("simulated GetByID failure")
}

func (e *ErrorGetByIDUserStore) GetAll(_ context.Context) ([]User, error) { return nil, nil }

// ErrorResultStore is a test double that always fails on SaveResult.
type ErrorResultStore struct{}

func NewErrorResultStore() *ErrorResultStore { return &ErrorResultStore{} }

func (e *ErrorResultStore) SaveResult(_ context.Context, _ Result) error {
	return fmt.Errorf("simulated result store failure")
}

func (e *ErrorResultStore) GetResult(_ context.Context, _ string) (*Result, error) {
	return nil, nil
}
