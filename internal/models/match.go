// Package models defines the shared domain types passed across the server,
// repository, and scoring layers — primarily Match, MatchSummary, and the
// associated event/result structs.
package models

import "time"

type Match struct {
	ID           string
	HomeTeam     string
	AwayTeam     string
	Kickoff      time.Time
	Status       string
	HomeScore    string
	AwayScore    string
	State        string
	DisplayClock string
	Venue        string
	HomeRecord   string
	AwayRecord   string
	HomeForm     string
	AwayForm     string
	HomeLogo     string
	AwayLogo     string
	Attendance   int
	Referee      string
	Events       []MatchEvent

	// LastPollAt is set by every /admin/poll-scores call. The 4am/12pm/6pm
	// refresh inspects it to detect dead polling chains (state=in but
	// LastPollAt is stale) and revives them. Zero value means never polled.
	LastPollAt time.Time
}

type MatchEvent struct {
	Clock   string   `json:"clock"`
	TypeID  string   `json:"typeID"`
	Team    string   `json:"team"`
	Players []string `json:"players"`
}

type MatchSummary struct {
	Attendance int
	HomeLogo   string
	AwayLogo   string
	Referee    string
	Events     []MatchEvent
}
