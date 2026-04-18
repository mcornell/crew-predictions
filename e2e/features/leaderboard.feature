Feature: Leaderboard

  Scenario: User sees their points after an exact score prediction
    Given I am logged in as "BlackAndGold@bsky.mock"
    And "BlackAndGold@bsky.mock" predicted 2-0 for match "match-scoring-1"
    And the final score for match "match-scoring-1" was 2-0
    When I visit the leaderboard
    Then I should see "BlackAndGold@bsky.mock" with 15 points
