@reset
Feature: Match detail page

  Scenario: Logged-in user sees all predictions for a completed match
    Given I am logged in as "CrewFan@bsky.mock"
    And the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status          | homeScore | awayScore |
      | m-done-1 | Columbus Crew | FC Dallas | STATUS_FULL_TIME | 2         | 1         |
    And "CrewFan@bsky.mock" predicted 2-1 for match "m-done-1"
    And "OtherFan@bsky.mock" predicted 1-1 for match "m-done-1"
    And the final score for match "m-done-1" was 2-1 with Columbus home
    When I visit the matches page
    And I click on the result card for match "m-done-1"
    Then I should be on the match detail page for "m-done-1"
    And I should see the match header with "Columbus Crew" vs "FC Dallas"
    And I should see "CrewFan@bsky.mock" in the predictions table
    And I should see "OtherFan@bsky.mock" in the predictions table
    And "CrewFan@bsky.mock" should have more points than "OtherFan@bsky.mock"

  Scenario: Result card links to match detail page
    Given I am not logged in
    And the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status          | homeScore | awayScore |
      | m-done-2 | Columbus Crew | LA Galaxy | STATUS_FULL_TIME | 3         | 0         |
    And the final score for match "m-done-2" was 3-0 with Columbus home
    When I visit the matches page
    Then the result card for match "m-done-2" should link to "/matches/m-done-2"

  Scenario: Upcoming match card does not link to match detail
    Given I am not logged in
    And the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status           |
      | m-sched  | Columbus Crew | FC Dallas | STATUS_SCHEDULED |
    When I visit the matches page
    Then the upcoming card for match "m-sched" should not have a detail link

  Scenario: No predictions made shows empty state
    Given I am not logged in
    And the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status          | homeScore | awayScore |
      | m-done-3 | Columbus Crew | FC Dallas | STATUS_FULL_TIME | 1         | 0         |
    And the final score for match "m-done-3" was 1-0 with Columbus home
    When I visit the match detail page for "m-done-3"
    Then I should see "No predictions were made for this match"
