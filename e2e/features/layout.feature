Feature: Page layout

  Scenario: Every page has a site header and proper structure
    Given I am not logged in
    When I visit the matches page
    Then the page title should be "Crew Predictions"
    And I should see a site header with "Crew Predictions"
