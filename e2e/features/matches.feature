Feature: Match listings

  Scenario: Unauthenticated user sees upcoming Columbus Crew matches
    Given I am not logged in
    When I visit the matches page
    Then I should see the "Upcoming Matches" heading
    And I should see at least one Columbus Crew match card
