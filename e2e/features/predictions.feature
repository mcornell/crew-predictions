Feature: Score predictions

  Scenario: Logged-in user submits a score prediction
    Given I am logged in as "BlackAndGold@bsky.mock"
    When I visit the matches page
    And I enter a home score of 3 and away score of 1 for the first match
    And I click "Lock In"
    Then I should see my prediction of "3 – 1" on the first match card
