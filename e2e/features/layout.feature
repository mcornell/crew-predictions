Feature: Page layout

  Scenario: Unknown route shows a not-found page
    Given I am not logged in
    When I visit an unknown page
    Then I should see a not-found message
    And I should see a link home

  Scenario: Every page has a site header and proper structure
    Given I am not logged in
    When I visit the matches page
    Then the page title should be "Crew Predictions"
    And I should see a site header with "Crew Predictions"
