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

	// ChainSeededFor records the Kickoff value the current Cloud Tasks
	// chain was seeded for. The refresh handler uses it as the dedup key:
	// if ChainSeededFor == Kickoff, a chain task has already been enqueued
	// for this match's current kickoff and we don't re-seed.
	ChainSeededFor time.Time

	// AbandonedAt is set when the 4h safety bailout in /admin/poll-scores
	// fires (now - kickoff > 4h). Diagnostic only — the chain ends without
	// it on normal terminal status; this field marks the bail-out case so
	// it's visible in match data when debugging.
	AbandonedAt time.Time
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
