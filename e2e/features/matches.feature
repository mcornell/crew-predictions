@reset
Feature: Match listings

  # ============================================================
  # Page-level structure
  # ============================================================

  Scenario: Matches page renders Crew matches and section headings
    Given the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status             | state |
      | m-pg-up  | Columbus Crew | FC Dallas | STATUS_SCHEDULED   |       |
      | m-pg-now | Columbus Crew | LA Galaxy | STATUS_IN_PROGRESS | in    |
    And I am not logged in
    When I visit the matches page
    Then I should see the "Now Playing" heading
    And I should see the "Upcoming" heading
    And I should see at least one Columbus Crew match card
    And the now playing card should appear before the upcoming card

  Scenario: Matches list orders cards by kickoff (earliest first)
    Given the following matches are seeded in order:
      | id      | homeTeam      | awayTeam  | status           | kickoffOffset |
      | m-late  | Columbus Crew | FC Dallas | STATUS_SCHEDULED | 48            |
      | m-early | Columbus Crew | LA Galaxy | STATUS_SCHEDULED | 24            |
    And I am not logged in
    When I visit the matches page
    Then match "m-early" should appear before match "m-late"

  Scenario: Admin refresh endpoint populates match cache from the fetcher
    Given the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status           |
      | m-cached | Columbus Crew | FC Dallas | STATUS_SCHEDULED |
    When the admin triggers a match refresh
    Then the matches API includes match "m-cached"

  # ============================================================
  # Card variant displays (consolidated — full content in one scenario,
  # using soft assertions so all failures are reported together)
  # ============================================================

  Scenario: Upcoming match card displays all expected content
    Given the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status           | homeRecord | awayRecord | homeForm | awayForm | venue                   |
      | m-up-all | Columbus Crew | FC Dallas | STATUS_SCHEDULED | 5-3-2      | 4-4-2      | WWWLL    | LWDWL    | ScottsMiracle-Gro Field |
    And I am not logged in
    When I visit the matches page
    Then I should see a countdown on the match card
    And the match card for "m-up-all" should show venue "ScottsMiracle-Gro Field"
    And the match card for "m-up-all" should show home record "5-3-2"
    And the match card for "m-up-all" should show home form "WWWLL"

  Scenario: Predicted upcoming match card still shows records and form
    Given the following matches are seeded:
      | id          | homeTeam      | awayTeam  | status           | homeRecord | awayRecord | homeForm | awayForm |
      | m-up-pred   | Columbus Crew | FC Dallas | STATUS_SCHEDULED | 5-3-2      | 4-4-2      | WWWLL    | LWDWL    |
    And I am logged in as "BlackAndGold@bsky.mock"
    And I have a seeded prediction of 2-1 for match "m-up-pred"
    When I visit the matches page
    Then the match card for "m-up-pred" should show home record "5-3-2"
    And the match card for "m-up-pred" should show home form "WWWLL"

  Scenario: Live match card displays all expected content
    Given the following matches are seeded:
      | id         | homeTeam      | awayTeam  | status             | state | homeScore | awayScore | homeRecord | awayRecord | homeForm | awayForm | venue                   |
      | m-live-all | Columbus Crew | FC Dallas | STATUS_IN_PROGRESS | in    | 2         | 1         | 5-3-2      | 4-4-2      | WWWLL    | LWDWL    | ScottsMiracle-Gro Field |
    And the following events are seeded for match "m-live-all":
      | clock | typeID       | team          | players                       |
      | 23'   | goal         | Columbus Crew | Hugo Picard                   |
      | 39'   | yellow-card  | FC Dallas     | Some Player                   |
      | 60'   | substitution | Columbus Crew | Sub In Player, Sub Off Player |
    And I am not logged in
    When I visit the matches page
    Then I should see a LIVE indicator on the match card
    And the now playing card should show score "2" to "1"
    And the now playing card for match "m-live-all" should show venue "ScottsMiracle-Gro Field"
    And the now playing card for match "m-live-all" should show home record "5-3-2"
    And the now playing card for match "m-live-all" should show home form "WWWLL"
    And the now playing card for match "m-live-all" should show event content
    And the now playing card for match "m-live-all" should show "Hugo Picard"
    And the now playing card for match "m-live-all" should show "Some Player"
    And the now playing card for match "m-live-all" should not show "Sub In Player"

  Scenario: Result card displays all expected content
    Given the following matches are seeded:
      | id          | homeTeam      | awayTeam  | status           | homeScore | awayScore | venue                   |
      | m-result-all| Columbus Crew | FC Dallas | STATUS_FULL_TIME | 2         | 1         | ScottsMiracle-Gro Field |
    And the final score for match "m-result-all" was 2-1 with Columbus home
    And I am logged in as "BlackAndGold@bsky.mock"
    And I have a seeded prediction of 1-0 for match "m-result-all"
    When I visit the matches page
    Then the result card for match "m-result-all" should show venue "ScottsMiracle-Gro Field"
    And the result card for match "m-result-all" should show my pick "1 – 0" below the score

  # ============================================================
  # State-driven behaviors (kept atomic — different conditions, not display)
  # ============================================================

  Scenario: Live match card hides events block when no events have occurred
    Given the following matches are seeded:
      | id      | homeTeam      | awayTeam  | status             | state | homeScore | awayScore |
      | m-noevt | Columbus Crew | FC Dallas | STATUS_IN_PROGRESS | in    | 0         | 0         |
    And I am not logged in
    When I visit the matches page
    Then the now playing card for match "m-noevt" should not show an events block

  Scenario: Predict button is absent for a match past kickoff
    Given the following matches are seeded in order:
      | id       | homeTeam      | awayTeam  | status           | kickoffOffset |
      | m-kicked | Columbus Crew | FC Dallas | STATUS_SCHEDULED | -1            |
    And I am not logged in
    When I visit the matches page
    Then I should not see a "Predict" button

  Scenario: Delayed match shows DELAYED badge and no Predict button
    Given the following matches are seeded:
      | id        | homeTeam      | awayTeam  | status         |
      | m-delayed | Columbus Crew | LA Galaxy | STATUS_DELAYED |
    And I am not logged in
    When I visit the matches page
    Then I should see a DELAYED indicator on the match card
    And I should not see a "Predict" button
