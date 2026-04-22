package repository

import (
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

type stubResultDoc struct {
	result Result
	err    error
}

func (s *stubResultDoc) DataTo(p interface{}) error {
	if s.err != nil {
		return s.err
	}
	*(p.(*Result)) = s.result
	return nil
}

func TestToPrediction_ReturnsMappedPrediction(t *testing.T) {
	doc := &stubDoc{prediction: Prediction{MatchID: "m1", UserID: "u1", HomeGoals: 3, AwayGoals: 1}}
	got, err := toPrediction(doc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.MatchID != "m1" || got.HomeGoals != 3 {
		t.Errorf("unexpected prediction: %+v", got)
	}
}

func TestToPrediction_ReturnsErrorWhenDataToFails(t *testing.T) {
	doc := &stubDoc{err: fmt.Errorf("bad data")}
	_, err := toPrediction(doc)
	if err == nil {
		t.Error("expected error when DataTo fails, got nil")
	}
}

func TestToResult_ReturnsMappedResult(t *testing.T) {
	doc := &stubResultDoc{result: Result{MatchID: "m1", HomeGoals: 2, AwayGoals: 1}}
	got, err := toResult(doc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.MatchID != "m1" || got.HomeGoals != 2 {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestToResult_ReturnsErrorWhenDataToFails(t *testing.T) {
	doc := &stubResultDoc{err: fmt.Errorf("bad data")}
	_, err := toResult(doc)
	if err == nil {
		t.Error("expected error when DataTo fails, got nil")
	}
}

func TestIsNotFound_TrueForNotFoundError(t *testing.T) {
	err := status.Error(codes.NotFound, "not found")
	if !isNotFound(err) {
		t.Error("expected true for NotFound gRPC error")
	}
}

func TestIsNotFound_FalseForOtherErrors(t *testing.T) {
	err := status.Error(codes.Internal, "server error")
	if isNotFound(err) {
		t.Error("expected false for non-NotFound gRPC error")
	}
}
