package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/mcornell/crew-predictions/internal/models"
)

type FirestoreMatchStore struct {
	client *firestore.Client
}

func NewFirestoreMatchStore(projectID string) (*FirestoreMatchStore, error) {
	client, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, err
	}
	return &FirestoreMatchStore{client: client}, nil
}

func (s *FirestoreMatchStore) SaveAll(matches []models.Match) error {
	ctx := context.Background()
	batch := s.client.BulkWriter(ctx)
	for _, m := range matches {
		ref := s.client.Collection("matches").Doc(m.ID)
		batch.Set(ref, map[string]any{
			"homeTeam":     m.HomeTeam,
			"awayTeam":     m.AwayTeam,
			"kickoff":      m.Kickoff,
			"status":       m.Status,
			"homeScore":    m.HomeScore,
			"awayScore":    m.AwayScore,
			"state":        m.State,
			"displayClock": m.DisplayClock,
			"venue":        m.Venue,
		})
	}
	batch.Flush()
	return nil
}

func (s *FirestoreMatchStore) GetAll() ([]models.Match, error) {
	ctx := context.Background()
	snaps, err := s.client.Collection("matches").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	out := make([]models.Match, 0, len(snaps))
	for _, snap := range snaps {
		m, err := toMatch(snap)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

func (s *FirestoreMatchStore) Reset() {}

func toMatch(snap *firestore.DocumentSnapshot) (models.Match, error) {
	var doc struct {
		HomeTeam     string    `firestore:"homeTeam"`
		AwayTeam     string    `firestore:"awayTeam"`
		Kickoff      time.Time `firestore:"kickoff"`
		Status       string    `firestore:"status"`
		HomeScore    string    `firestore:"homeScore"`
		AwayScore    string    `firestore:"awayScore"`
		State        string    `firestore:"state"`
		DisplayClock string    `firestore:"displayClock"`
		Venue        string    `firestore:"venue"`
	}
	if err := snap.DataTo(&doc); err != nil {
		return models.Match{}, err
	}
	return models.Match{
		ID:           snap.Ref.ID,
		HomeTeam:     doc.HomeTeam,
		AwayTeam:     doc.AwayTeam,
		Kickoff:      doc.Kickoff,
		Status:       doc.Status,
		HomeScore:    doc.HomeScore,
		AwayScore:    doc.AwayScore,
		State:        doc.State,
		DisplayClock: doc.DisplayClock,
		Venue:        doc.Venue,
	}, nil
}
