@smoke
Feature: API health checks

  Scenario: Matches API returns valid JSON
    Then the API at "/api/matches" returns a JSON array

  Scenario: Leaderboard API returns valid JSON
    Then the API at "/api/leaderboard" returns a JSON array
