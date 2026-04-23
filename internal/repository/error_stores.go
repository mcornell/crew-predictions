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
