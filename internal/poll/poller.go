package poll

import (
	"context"
	"strconv"

	"github.com/mcornell/crew-predictions/internal/models"
	"github.com/mcornell/crew-predictions/internal/repository"
)

var terminalStatuses = map[string]bool{
	"STATUS_FULL_TIME":  true,
	"STATUS_FINAL_AET":  true,
	"STATUS_FINAL_PEN":  true,
}

func PollOnce(ctx context.Context, matchStore repository.MatchStore, resultStore repository.ResultStore, fetcher func() ([]models.Match, error)) error {
	matches, err := fetcher()
	if err != nil {
		return err
	}

	if err := matchStore.SaveAll(matches); err != nil {
		return err
	}

	for _, m := range matches {
		if !terminalStatuses[m.Status] {
			continue
		}
		home, err := strconv.Atoi(m.HomeScore)
		if err != nil {
			continue
		}
		away, err := strconv.Atoi(m.AwayScore)
		if err != nil {
			continue
		}
		if err := resultStore.SaveResult(ctx, repository.Result{
			MatchID:   m.ID,
			HomeTeam:  m.HomeTeam,
			AwayTeam:  m.AwayTeam,
			HomeGoals: home,
			AwayGoals: away,
		}); err != nil {
			return err
		}
	}
	return nil
}
