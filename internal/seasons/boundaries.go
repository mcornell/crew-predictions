package seasons

import "time"

type SeasonDef struct {
	ID    string
	Name  string
	Start time.Time
	End   time.Time
}

var seasonTable = []SeasonDef{
	{
		ID:    "2026",
		Name:  "2026 Season",
		Start: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2027, 1, 10, 0, 0, 0, 0, time.UTC),
	},
	{
		ID:    "2027-sprint",
		Name:  "2027 Sprint Season",
		Start: time.Date(2027, 1, 10, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2027, 6, 20, 0, 0, 0, 0, time.UTC),
	},
	{
		ID:    "2027-28",
		Name:  "2027-28 Season",
		Start: time.Date(2027, 6, 20, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2028, 6, 20, 0, 0, 0, 0, time.UTC),
	},
	{
		ID:    "2028-29",
		Name:  "2028-29 Season",
		Start: time.Date(2028, 6, 20, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2029, 6, 20, 0, 0, 0, 0, time.UTC),
	},
}

func AllSeasons() []SeasonDef {
	out := make([]SeasonDef, len(seasonTable))
	copy(out, seasonTable)
	return out
}

func SeasonByID(id string) (SeasonDef, bool) {
	for _, s := range seasonTable {
		if s.ID == id {
			return s, true
		}
	}
	return SeasonDef{}, false
}

// SeasonForDate returns the season that contains t. Falls back to the last
// known season if t is beyond all defined boundaries.
func SeasonForDate(t time.Time) SeasonDef {
	for i := len(seasonTable) - 1; i >= 0; i-- {
		if !t.Before(seasonTable[i].Start) {
			return seasonTable[i]
		}
	}
	return seasonTable[0]
}

// MaybeCloseSeason returns the season IDs to close and open, and whether the
// close should proceed. Returns false if now is before the boundary, or if the
// season snapshot already exists (idempotency).
func MaybeCloseSeason(now time.Time, activeID string, seasonExists func(string) bool) (closeID, openID string, should bool) {
	active, ok := SeasonByID(activeID)
	if !ok {
		return "", "", false
	}
	if now.Before(active.End) {
		return "", "", false
	}
	if seasonExists(activeID) {
		return "", "", false
	}
	next := SeasonForDate(active.End)
	return activeID, next.ID, true
}
