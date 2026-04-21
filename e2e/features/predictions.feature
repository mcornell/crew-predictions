Feature: Score predictions

  Background:
    Given the following matches are seeded:
      | id           | homeTeam         | awayTeam      | status           |
      | match-pred-1 | Columbus Crew    | LA Galaxy     | STATUS_SCHEDULED |
      | match-past-1 | Portland Timbers | Columbus Crew | STATUS_FULL_TIME |

  Scenario: Logged-out user clicking Predict is redirected to sign in
    Given I am not logged in
    When I visit the matches page
    And I click "Predict"
    Then I should be on the login page

  Scenario: Logged-in user submits a score prediction
    Given I am logged in as "BlackAndGold@bsky.mock"
    When I visit the matches page
    And I enter a home score of 3 and away score of 1 for the first match
    And I click "Predict"
    Then I should see my prediction of "3" on the first match card

  Scenario: Prediction is rejected after kickoff
    Given I am logged in as "BlackAndGold@bsky.mock"
    When I submit a prediction via API for match "match-past-1"
    Then the server should reject it with 403
