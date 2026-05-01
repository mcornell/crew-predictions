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
	Attendance   int
	Events       []MatchEvent
}

type MatchEvent struct {
	Clock   string   `json:"clock"`
	TypeID  string   `json:"typeID"`
	Team    string   `json:"team"`
	Players []string `json:"players"`
}

type MatchSummary struct {
	Attendance int
	Events     []MatchEvent
}
