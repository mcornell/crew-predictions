package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreResultStore struct {
	client *firestore.Client
}

func NewFirestoreResultStore(ctx context.Context, projectID string) (*FirestoreResultStore, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &FirestoreResultStore{client: client}, nil
}

func (s *FirestoreResultStore) SaveResult(ctx context.Context, r Result) error {
	_, err := s.client.Collection("results").Doc(r.MatchID).Set(ctx, map[string]any{
		"MatchID":   r.MatchID,
		"HomeTeam":  r.HomeTeam,
		"AwayTeam":  r.AwayTeam,
		"HomeGoals": r.HomeGoals,
		"AwayGoals": r.AwayGoals,
	})
	return err
}

func (s *FirestoreResultStore) GetResult(ctx context.Context, matchID string) (*Result, error) {
	snap, err := s.client.Collection("results").Doc(matchID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, err
	}
	var r Result
	if err := snap.DataTo(&r); err != nil {
		return nil, err
	}
	return &r, nil
}
