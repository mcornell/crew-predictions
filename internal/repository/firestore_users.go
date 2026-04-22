package repository

import (
	"context"

	"cloud.google.com/go/firestore"
)

type FirestoreUserStore struct {
	client *firestore.Client
}

func NewFirestoreUserStore(ctx context.Context, projectID string) (*FirestoreUserStore, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &FirestoreUserStore{client: client}, nil
}

func (s *FirestoreUserStore) Upsert(ctx context.Context, u User) error {
	data := map[string]any{"handle": u.Handle}
	if u.Provider != "" {
		data["provider"] = u.Provider
	}
	if u.Location != "" {
		data["location"] = u.Location
	}
	_, err := s.client.Collection("users").Doc(u.UserID).Set(ctx, data, firestore.MergeAll)
	return err
}

func (s *FirestoreUserStore) GetByID(ctx context.Context, userID string) (*User, error) {
	snap, err := s.client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return toUser(userID, snap)
}

func (s *FirestoreUserStore) GetAll(ctx context.Context) ([]User, error) {
	snapshots, err := s.client.Collection("users").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	users := make([]User, 0, len(snapshots))
	for _, snap := range snapshots {
		u, err := toUser(snap.Ref.ID, snap)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	return users, nil
}

func toUser(userID string, snap *firestore.DocumentSnapshot) (*User, error) {
	var doc struct {
		Handle   string `firestore:"handle"`
		Provider string `firestore:"provider"`
		Location string `firestore:"location"`
	}
	if err := snap.DataTo(&doc); err != nil {
		return nil, err
	}
	return &User{UserID: userID, Handle: doc.Handle, Provider: doc.Provider, Location: doc.Location}, nil
}
