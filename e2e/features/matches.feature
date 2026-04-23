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

  Scenario: Upcoming match card shows a countdown to kickoff
    Given the following matches are seeded:
      | id      | homeTeam      | awayTeam  | status           |
      | m-cnt-1 | Columbus Crew | FC Dallas | STATUS_SCHEDULED |
    And I am not logged in
    When I visit the matches page
    Then I should see a countdown on the match card

  Scenario: User sees LIVE indicator on an in-progress match
    Given the following matches are seeded:
      | id     | homeTeam      | awayTeam  | status           | state |
      | m-live | Columbus Crew | FC Dallas | STATUS_SCHEDULED | in    |
    And I am not logged in
    When I visit the matches page
    Then I should see a LIVE indicator on the match card

  Scenario: Delayed match shows DELAYED badge and no Predict button
    Given the following matches are seeded:
      | id        | homeTeam      | awayTeam  | status         |
      | m-delayed | Columbus Crew | LA Galaxy | STATUS_DELAYED |
    And I am not logged in
    When I visit the matches page
    Then I should see a DELAYED indicator on the match card
    And I should not see a "Predict" button

  Scenario: Matches display in kickoff order earliest first
    Given the following matches are seeded in order:
      | id      | homeTeam      | awayTeam    | status           | kickoffOffset |
      | m-late  | Columbus Crew | FC Dallas   | STATUS_SCHEDULED | 48            |
      | m-early | Columbus Crew | LA Galaxy   | STATUS_SCHEDULED | 24            |
    And I am not logged in
    When I visit the matches page
    Then match "m-early" should appear before match "m-late"

  Scenario: Predict button is absent for a match past kickoff
    Given the following matches are seeded in order:
      | id       | homeTeam      | awayTeam  | status           | kickoffOffset |
      | m-kicked | Columbus Crew | FC Dallas | STATUS_SCHEDULED | -1            |
    And I am not logged in
    When I visit the matches page
    Then I should not see a "Predict" button
