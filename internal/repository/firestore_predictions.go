package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	doc := s.client.Collection("predictions").Doc(p.MatchID + "|" + p.Handle)
	_, err := doc.Set(ctx, map[string]any{
		"MatchID":   p.MatchID,
		"Handle":    p.Handle,
		"HomeGoals": p.HomeGoals,
		"AwayGoals": p.AwayGoals,
	})
	return err
}

func (s *FirestorePredictionStore) GetAll(ctx context.Context) ([]Prediction, error) {
	docs, err := s.client.Collection("predictions").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	all := make([]Prediction, 0, len(docs))
	for _, doc := range docs {
		var p Prediction
		if err := doc.DataTo(&p); err != nil {
			return nil, err
		}
		all = append(all, p)
	}
	return all, nil
}

func (s *FirestorePredictionStore) GetByMatchAndHandle(ctx context.Context, matchID, handle string) (*Prediction, error) {
	doc := s.client.Collection("predictions").Doc(matchID + "|" + handle)
	snap, err := doc.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, err
	}
	var p Prediction
	if err := snap.DataTo(&p); err != nil {
		return nil, err
	}
	return &p, nil
}
