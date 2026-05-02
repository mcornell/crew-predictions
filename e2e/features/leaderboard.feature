@reset
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
    Then I should see "ColumbusNordecke@bsky.mock" with 2 Upper 90 Club points

  Scenario: User sees Grouchy points after a correct category prediction
    Given I am logged in as "GrouchyFan@bsky.mock"
    And "GrouchyFan@bsky.mock" predicted 3-1 for match "match-scoring-3"
    And the final score for match "match-scoring-3" was 4-0 with Columbus away
    When I visit the leaderboard
    Then I should see "GrouchyFan@bsky.mock" with 1 Grouchy points

  Scenario: Leaderboard nav goes directly to current season; selector lives on the leaderboard page
    When an admin closes the current season
    And I am on the matches page
    Then no season selector should be visible
    When I click the Leaderboard nav link
    Then I should be on the leaderboard page
    And a season selector should be visible

  Scenario: Season dropdown shows past seasons and hides future seasons
    When an admin closes the current season
    And I visit the leaderboard
    And I open the leaderboard dropdown
    Then the season dropdown includes "2026 Season"
    And the season dropdown does not include "2027-28 Season"
    And the season dropdown does not include "2028-29 Season"
