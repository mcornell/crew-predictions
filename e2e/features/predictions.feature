@reset
Feature: Score predictions

  Background:
    Given the following matches are seeded:
      | id           | homeTeam         | awayTeam      | status           |
      | match-pred-1 | Columbus Crew    | LA Galaxy     | STATUS_SCHEDULED |
      | match-past-1 | Portland Timbers | Columbus Crew | STATUS_FULL_TIME |

  Scenario: Logged-out user sees score inputs and an enabled Predict button
    Given I am not logged in
    When I visit the matches page
    Then I should see a "Predict" button

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

  Scenario: Guest user can enter a prediction without signing in
    Given I am not logged in
    When I visit the matches page
    And I enter a home score of 2 and away score of 0 for the first match
    And I click "Predict"
    Then I should see my prediction of "2" on the first match card

  Scenario: Guest prediction persists after page reload
    Given I am not logged in
    When I visit the matches page
    And I enter a home score of 3 and away score of 0 for the first match
    And I click "Predict"
    And I reload the page
    Then I should see my prediction of "3" on the first match card

  Scenario: Guest user sees a sign-in nudge after predicting
    Given I am not logged in
    When I visit the matches page
    And I enter a home score of 2 and away score of 0 for the first match
    And I click "Predict"
    Then I should see a sign-in nudge

  Scenario: Logged-in user can unlock and re-submit a prediction
    Given I am logged in as "BlackAndGold@bsky.mock"
    When I visit the matches page
    And I enter a home score of 2 and away score of 1 for the first match
    And I click "Predict"
    And I click "Unlock"
    Then I should see an enabled "Predict" button
    And the first match score inputs should show 2 and 1
    When I enter a home score of 3 and away score of 0 for the first match
    And I click "Predict"
    Then I should see my prediction of "3" on the first match card
