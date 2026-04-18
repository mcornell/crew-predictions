Feature: Leaderboard

  Scenario: User sees Aces Radio points after an exact score prediction
    Given I am logged in as "BlackAndGold@bsky.mock"
    And "BlackAndGold@bsky.mock" predicted 2-0 for match "match-scoring-1"
    And the final score for match "match-scoring-1" was 2-0 with Columbus away
    When I visit the leaderboard
    Then I should see "BlackAndGold@bsky.mock" with 15 Aces Radio points

  Scenario: User sees Upper90Club points after a correct winner prediction
    Given I am logged in as "ColumbusNordecke@bsky.mock"
    And "ColumbusNordecke@bsky.mock" predicted 1-0 for match "match-scoring-2"
    And the final score for match "match-scoring-2" was 3-0 with Columbus away
    When I visit the leaderboard
    Then I should see "ColumbusNordecke@bsky.mock" with 2 Upper90Club points
