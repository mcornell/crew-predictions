Feature: Season history

  @reset
  Scenario: Historical leaderboard shows frozen standings for a closed season
    Given season "2026" has been archived with "HistoryFan" at 15 Aces Radio points
    When I visit the historical leaderboard for season "2026"
    Then I should see "HistoryFan" with 15 Aces Radio points
    And I should see a season selector on the leaderboard page
