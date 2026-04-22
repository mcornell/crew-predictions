@reset
Feature: Match listings

  Scenario: Unauthenticated user sees upcoming Columbus Crew matches
    Given the following matches are seeded:
      | id      | homeTeam      | awayTeam  | status           |
      | m-clb-1 | Columbus Crew | FC Dallas | STATUS_SCHEDULED |
    And I am not logged in
    When I visit the matches page
    Then I should see the "Upcoming" heading
    And I should see at least one Columbus Crew match card

  Scenario: Admin refresh endpoint populates match cache from the fetcher
    Given the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status           |
      | m-cached | Columbus Crew | FC Dallas | STATUS_SCHEDULED |
    When the admin triggers a match refresh
    Then the matches API includes match "m-cached"
