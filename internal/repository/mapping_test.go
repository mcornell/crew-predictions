package repository

import (
	"fmt"
	"testing"
)

type stubDoc struct {
	prediction Prediction
	err        error
}

func (s *stubDoc) DataTo(p interface{}) error {
	if s.err != nil {
		return s.err
	}
	*(p.(*Prediction)) = s.prediction
	return nil
}

func TestToPredictions_ReturnsMappedSlice(t *testing.T) {
	docs := []dataMapper{
		&stubDoc{prediction: Prediction{MatchID: "m1", UserID: "u1", HomeGoals: 2, AwayGoals: 0}},
		&stubDoc{prediction: Prediction{MatchID: "m2", UserID: "u2", HomeGoals: 1, AwayGoals: 1}},
	}

	got, err := toPredictions(docs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 predictions, got %d", len(got))
	}
	if got[0].MatchID != "m1" || got[1].MatchID != "m2" {
		t.Errorf("unexpected predictions: %+v", got)
	}
}

func TestToPredictions_ReturnsEmptySliceForNoDocs(t *testing.T) {
	got, err := toPredictions(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %+v", got)
	}
}

func TestToPredictions_ReturnsErrorWhenDataToFails(t *testing.T) {
	docs := []dataMapper{
		&stubDoc{prediction: Prediction{MatchID: "m1"}},
		&stubDoc{err: fmt.Errorf("corrupt document")},
	}

	_, err := toPredictions(docs)
	if err == nil {
		t.Error("expected error when DataTo fails, got nil")
	}
}
