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

  Scenario: Grouchy column appears in match detail predictions table
    Given I am logged in as "GrouchyFan@bsky.mock"
    And the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status          | homeScore | awayScore |
      | m-done-g | Columbus Crew | FC Dallas | STATUS_FULL_TIME | 2         | 0         |
    And "GrouchyFan@bsky.mock" predicted 3-0 for match "m-done-g"
    And the final score for match "m-done-g" was 2-0 with Columbus home
    When I visit the match detail page for "m-done-g"
    Then I should see the Grouchy column header in the predictions table
    And "GrouchyFan@bsky.mock" should have 1 Grouchy point in the predictions table

  Scenario: No predictions made shows empty state
    Given I am not logged in
    And the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status          | homeScore | awayScore |
      | m-done-3 | Columbus Crew | FC Dallas | STATUS_FULL_TIME | 1         | 0         |
    And the final score for match "m-done-3" was 1-0 with Columbus home
    When I visit the match detail page for "m-done-3"
    Then I should see "No predictions were made for this match"

  Scenario: Live match detail shows LIVE indicator and projected scores
    Given I am logged in as "CrewFan@bsky.mock"
    And the following matches are seeded:
      | id      | homeTeam      | awayTeam          | status             | state | homeScore | awayScore |
      | m-live  | Columbus Crew | Philadelphia Union | STATUS_IN_PROGRESS | in    | 2         | 0         |
    And "CrewFan@bsky.mock" predicted 2-0 for match "m-live"
    And "OtherFan@bsky.mock" predicted 1-1 for match "m-live"
    When I visit the match detail page for "m-live"
    Then I should see the LIVE indicator on the match detail page
    And the match detail header should show score "2" to "0"
    And the projected points label should be visible
    And "CrewFan@bsky.mock" should have projected points greater than "OtherFan@bsky.mock"

  Scenario: Live match card on main page links to match detail
    Given I am not logged in
    And the following matches are seeded:
      | id      | homeTeam      | awayTeam          | status             | state |
      | m-live2 | Columbus Crew | Philadelphia Union | STATUS_IN_PROGRESS | in    |
    When I visit the matches page
    Then the now playing card for match "m-live2" should link to "/matches/m-live2"

  Scenario: Match detail page shows venue name
    Given I am not logged in
    And the following matches are seeded:
      | id      | homeTeam      | awayTeam  | status           | homeScore | awayScore | venue                   |
      | m-ven-4 | Columbus Crew | FC Dallas | STATUS_FULL_TIME | 2         | 1         | ScottsMiracle-Gro Field |
    And the final score for match "m-ven-4" was 2-1 with Columbus home
    When I visit the match detail page for "m-ven-4"
    Then I should see the venue "ScottsMiracle-Gro Field" on the match detail page

  Scenario: Match detail page shows team records and form
    Given I am not logged in
    And the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status           | homeRecord | awayRecord | homeForm | awayForm |
      | m-form-3 | Columbus Crew | FC Dallas | STATUS_SCHEDULED | 5-3-2      | 4-4-2      | WWWLL    | LWDWL    |
    When I visit the match detail page for "m-form-3"
    Then I should see home record "5-3-2" on the match detail page
    And I should see home form "WWWLL" on the match detail page

  Scenario: Match detail page links to ESPN match page
    Given I am not logged in
    And the following matches are seeded:
      | id      | homeTeam      | awayTeam  | status           | homeScore | awayScore |
      | m-espn  | Columbus Crew | FC Dallas | STATUS_FULL_TIME | 2         | 1         |
    And the final score for match "m-espn" was 2-1 with Columbus home
    When I visit the match detail page for "m-espn"
    Then I should see an ESPN link for match "m-espn"

  Scenario: Match detail page lazy-fetches attendance from ESPN for a known completed match
    Given I am not logged in
    And the following matches are seeded in order:
      | id     | homeTeam      | awayTeam          | kickoffOffset | status           | state | homeScore | awayScore |
      | 761573 | Columbus Crew | Philadelphia Union | -120          | STATUS_FULL_TIME | post  | 2         | 0         |
    When I visit the match detail page for "761573"
    Then I should see the attendance "19,903" on the match detail page

  Scenario: Match detail page shows event timeline for a completed match
    Given I am not logged in
    And the following matches are seeded in order:
      | id     | homeTeam      | awayTeam          | kickoffOffset | status           | state | homeScore | awayScore |
      | 761573 | Columbus Crew | Philadelphia Union | -120          | STATUS_FULL_TIME | post  | 2         | 0         |
    When I visit the match detail page for "761573"
    Then I should see the event timeline on the match detail page
    And I should see at least one event in the timeline

  Scenario: Match detail page shows team logos for a completed match
    Given I am not logged in
    And the following matches are seeded in order:
      | id     | homeTeam      | awayTeam          | kickoffOffset | status           | state | homeScore | awayScore |
      | 761573 | Columbus Crew | Philadelphia Union | -120          | STATUS_FULL_TIME | post  | 2         | 0         |
    When I visit the match detail page for "761573"
    Then I should see the home team logo on the match detail page
    And I should see the away team logo on the match detail page

  Scenario: Match detail page shows the referee for a completed match
    Given I am not logged in
    And the following matches are seeded in order:
      | id     | homeTeam   | awayTeam      | kickoffOffset | status           | state | homeScore | awayScore |
      | 761499 | Toronto FC | Columbus Crew | -120          | STATUS_FULL_TIME | post  | 0         | 2         |
    When I visit the match detail page for "761499"
    Then I should see the referee on the match detail page
