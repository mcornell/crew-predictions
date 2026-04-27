@smoke
Feature: API health checks

  Scenario: Matches API returns valid JSON
    Then the API at "/api/matches" returns a JSON object with key "matches"

  Scenario: Leaderboard API returns valid JSON
    Then the API at "/api/leaderboard" returns a JSON object with key "entries"
