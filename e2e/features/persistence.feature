Feature: Prediction persistence

  Scenario: Submitted prediction survives a page reload
    Given I am logged in as "ColumbusNordecke@bsky.mock"
    When I visit the matches page
    And I enter a home score of 3 and away score of 1 for the first match
    And I click "Lock In"
    And I revisit the matches page
    Then I should see my prediction of "3 – 1" on the first match card
