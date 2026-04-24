package recalculator

import (
	"context"
	"fmt"

	"github.com/mcornell/crew-predictions/internal/repository"
	"github.com/mcornell/crew-predictions/internal/scoring"
)

func Recalculate(ctx context.Context, predictions repository.PredictionStore, results repository.ResultStore, users repository.UserStore, targetTeam string) error {
	allPredictions, err := predictions.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("recalculate: get predictions: %w", err)
	}

	allUsers, err := users.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("recalculate: get users: %w", err)
	}

	// Index predictions by userID.
	predsByUser := make(map[string][]repository.Prediction, len(allUsers))
	for _, p := range allPredictions {
		predsByUser[p.UserID] = append(predsByUser[p.UserID], p)
	}

	// Cache results by matchID to avoid redundant lookups.
	resultCache := map[string]*repository.Result{}
	for _, p := range allPredictions {
		if _, seen := resultCache[p.MatchID]; seen {
			continue
		}
		r, err := results.GetResult(ctx, p.MatchID)
		if err != nil {
			return fmt.Errorf("recalculate: get result %s: %w", p.MatchID, err)
		}
		resultCache[p.MatchID] = r
	}

	// Build a set of known userIDs for fast lookup.
	knownUsers := make(map[string]repository.User, len(allUsers))
	for _, u := range allUsers {
		knownUsers[u.UserID] = u
	}

	// Also include users who have predictions but no UserStore entry yet.
	// Upsert their profile so handle is available for leaderboard/match detail display.
	for userID, preds := range predsByUser {
		if _, exists := knownUsers[userID]; !exists {
			u := repository.User{UserID: userID, Handle: preds[0].Handle}
			knownUsers[userID] = u
			if err := users.Upsert(ctx, u); err != nil {
				return fmt.Errorf("recalculate: upsert new user %s: %w", userID, err)
			}
		}
	}

	for _, u := range knownUsers {
		var acesTotal, u90Total, grouchyTotal, predCount int
		for _, p := range predsByUser[u.UserID] {
			predCount++
			r := resultCache[p.MatchID]
			if r == nil {
				continue
			}
			pred := scoring.Prediction{Home: p.HomeGoals, Away: p.AwayGoals}
			res := scoring.Result{Home: r.HomeGoals, Away: r.AwayGoals}
			targetIsHome := r.HomeTeam == targetTeam
			acesTotal += scoring.AcesRadio(res, pred)
			u90Total += scoring.Upper90Club(res, pred, targetIsHome)
			grouchyTotal += scoring.Grouchy(res, pred, targetIsHome)
		}
		if err := users.UpdateScores(ctx, u.UserID, predCount, acesTotal, u90Total, grouchyTotal); err != nil {
			return fmt.Errorf("recalculate: update scores %s: %w", u.UserID, err)
		}
	}

	return nil
}
