@reset
Feature: Score polling

  Scenario: Score poll records result for a finished match
    Given the following matches are seeded:
      | id     | homeTeam      | awayTeam  | status           | state | homeScore | awayScore |
      | m-done | Columbus Crew | FC Dallas | STATUS_FULL_TIME | post  | 3         | 1         |
    And "alice@bsky.mock" predicted 3-1 for match "m-done"
    When the admin triggers a score poll
    And I visit the leaderboard
    Then I should see "alice@bsky.mock" with 15 Aces Radio points
