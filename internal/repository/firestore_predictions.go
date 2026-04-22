package repository

import (
	"context"

	"cloud.google.com/go/firestore"
)

type FirestorePredictionStore struct {
	client *firestore.Client
}

func NewFirestorePredictionStore(ctx context.Context, projectID string) (*FirestorePredictionStore, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &FirestorePredictionStore{client: client}, nil
}

func (s *FirestorePredictionStore) Save(ctx context.Context, p Prediction) error {
	doc := s.client.Collection("predictions").Doc(p.MatchID + "|" + p.UserID)
	_, err := doc.Set(ctx, map[string]any{
		"MatchID":   p.MatchID,
		"UserID":    p.UserID,
		"Handle":    p.Handle,
		"HomeGoals": p.HomeGoals,
		"AwayGoals": p.AwayGoals,
	})
	return err
}

func (s *FirestorePredictionStore) GetAll(ctx context.Context) ([]Prediction, error) {
	snapshots, err := s.client.Collection("predictions").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	docs := make([]dataMapper, len(snapshots))
	for i, d := range snapshots {
		docs[i] = d
	}
	return toPredictions(docs)
}

func (s *FirestorePredictionStore) GetByMatchAndUser(ctx context.Context, matchID, userID string) (*Prediction, error) {
	snap, err := s.client.Collection("predictions").Doc(matchID + "|" + userID).Get(ctx)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return toPrediction(snap)
}
