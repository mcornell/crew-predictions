Feature: Season history

  @reset
  Scenario: Historical leaderboard shows frozen standings for a closed season
    Given season "2026" has been archived with "HistoryFan" at 15 Aces Radio points
    When I visit the historical leaderboard for season "2026"
    Then I should see "HistoryFan" with 15 Aces Radio points
    And I should see a season selector on the leaderboard page

  @reset
  Scenario: Closing a season archives standings and resets current leaderboard
    Given "ArchiveFan" predicted 2-0 for match "match-archive-1"
    And the final score for match "match-archive-1" was 2-0 with Columbus away
    When an admin closes the current season
    Then visiting the historical leaderboard for season "2026" shows "ArchiveFan" with 15 Aces Radio points
    And the current leaderboard shows "ArchiveFan" with 0 Aces Radio points
