@reset
Feature: Match detail page

  # ============================================================
  # Card-to-detail navigation (which cards link to detail, which don't)
  # ============================================================

  Scenario: Match detail navigation from listings
    Given I am not logged in
    And the following matches are seeded:
      | id        | homeTeam      | awayTeam          | status             | state | homeScore | awayScore |
      | m-nav-up  | Columbus Crew | FC Dallas         | STATUS_SCHEDULED   |       |           |           |
      | m-nav-now | Columbus Crew | Philadelphia Union | STATUS_IN_PROGRESS | in    |           |           |
      | m-nav-end | Columbus Crew | LA Galaxy         | STATUS_FULL_TIME   |       | 3         | 0         |
    And the final score for match "m-nav-end" was 3-0 with Columbus home
    When I visit the matches page
    Then the result card for match "m-nav-end" should link to "/matches/m-nav-end"
    And the now playing card for match "m-nav-now" should link to "/matches/m-nav-now"
    And the upcoming card for match "m-nav-up" should not have a detail link

  # ============================================================
  # Completed match detail — non-ESPN-dependent display
  # ============================================================

  Scenario: Completed match detail page displays full match info
    Given I am logged in as "CrewFan@bsky.mock"
    And the following matches are seeded:
      | id      | homeTeam      | awayTeam  | status           | homeScore | awayScore | venue                   | homeRecord | awayRecord | homeForm | awayForm |
      | m-full  | Columbus Crew | FC Dallas | STATUS_FULL_TIME | 2         | 1         | ScottsMiracle-Gro Field | 5-3-2      | 4-4-2      | WWWLL    | LWDWL    |
    And "CrewFan@bsky.mock" predicted 2-1 for match "m-full"
    And "OtherFan@bsky.mock" predicted 1-1 for match "m-full"
    And the final score for match "m-full" was 2-1 with Columbus home
    When I visit the match detail page for "m-full"
    Then I should see the match header with "Columbus Crew" vs "FC Dallas"
    And I should see the venue "ScottsMiracle-Gro Field" on the match detail page
    And I should see home record "5-3-2" on the match detail page
    And I should see home form "WWWLL" on the match detail page
    And I should see an ESPN link for match "m-full"
    And I should see the Grouchy column header in the predictions table
    And I should see "CrewFan@bsky.mock" in the predictions table
    And I should see "OtherFan@bsky.mock" in the predictions table
    And "CrewFan@bsky.mock" should have more points than "OtherFan@bsky.mock"

  Scenario: Match detail page shows empty state when no predictions were made
    Given I am not logged in
    And the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status          | homeScore | awayScore |
      | m-empty  | Columbus Crew | FC Dallas | STATUS_FULL_TIME | 1         | 0         |
    And the final score for match "m-empty" was 1-0 with Columbus home
    When I visit the match detail page for "m-empty"
    Then I should see "No predictions were made for this match"

  # ============================================================
  # Live match detail (different state, kept atomic)
  # ============================================================

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

  # ============================================================
  # ESPN /summary enrichment — uses real ESPN, ~15s timeouts
  # Consolidated to one scenario per fixture match to amortize the
  # round-trip cost across multiple assertions.
  # ============================================================

  Scenario: Match detail page enriches a completed match from ESPN data
    Given I am not logged in
    And the following matches are seeded in order:
      | id     | homeTeam      | awayTeam           | kickoffOffset | status           | state | homeScore | awayScore |
      | 761573 | Columbus Crew | Philadelphia Union | -120          | STATUS_FULL_TIME | post  | 2         | 0         |
    When I visit the match detail page for "761573"
    Then I should see the attendance "19,903" on the match detail page
    And I should see the event timeline on the match detail page
    And I should see at least one event in the timeline
    And I should see the home team logo on the match detail page
    And I should see the away team logo on the match detail page

  Scenario: Match detail page shows the referee when ESPN provides one
    Given I am not logged in
    And the following matches are seeded in order:
      | id     | homeTeam   | awayTeam      | kickoffOffset | status           | state | homeScore | awayScore |
      | 761499 | Toronto FC | Columbus Crew | -120          | STATUS_FULL_TIME | post  | 0         | 2         |
    When I visit the match detail page for "761499"
    Then I should see the referee on the match detail page
