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

  Scenario: Live match appears in Now Playing section above Upcoming
    Given the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status             | state |
      | m-now    | Columbus Crew | FC Dallas | STATUS_IN_PROGRESS | in    |
      | m-sched2 | Columbus Crew | LA Galaxy | STATUS_SCHEDULED   |       |
    And I am not logged in
    When I visit the matches page
    Then I should see the "Now Playing" heading
    And I should see the "Upcoming" heading
    And the now playing card should appear before the upcoming card

  Scenario: Live match shows current score not dashes
    Given the following matches are seeded:
      | id      | homeTeam      | awayTeam  | status             | state | homeScore | awayScore |
      | m-score | Columbus Crew | FC Dallas | STATUS_IN_PROGRESS | in    | 2         | 1         |
    And I am not logged in
    When I visit the matches page
    Then the now playing card should show score "2" to "1"

  Scenario: Predict button is absent for a match past kickoff
    Given the following matches are seeded in order:
      | id       | homeTeam      | awayTeam  | status           | kickoffOffset |
      | m-kicked | Columbus Crew | FC Dallas | STATUS_SCHEDULED | -1            |
    And I am not logged in
    When I visit the matches page
    Then I should not see a "Predict" button

  Scenario: Upcoming match card shows venue name
    Given the following matches are seeded:
      | id      | homeTeam      | awayTeam  | status           | venue                   |
      | m-ven-1 | Columbus Crew | FC Dallas | STATUS_SCHEDULED | ScottsMiracle-Gro Field |
    And I am not logged in
    When I visit the matches page
    Then the match card for "m-ven-1" should show venue "ScottsMiracle-Gro Field"

  Scenario: Live match card shows venue name
    Given the following matches are seeded:
      | id      | homeTeam      | awayTeam  | status             | state | venue                   |
      | m-ven-2 | Columbus Crew | FC Dallas | STATUS_IN_PROGRESS | in    | ScottsMiracle-Gro Field |
    And I am not logged in
    When I visit the matches page
    Then the now playing card for match "m-ven-2" should show venue "ScottsMiracle-Gro Field"

  Scenario: Live match card shows goals and cards inline
    Given the following matches are seeded:
      | id     | homeTeam      | awayTeam  | status             | state | homeScore | awayScore |
      | m-evts | Columbus Crew | FC Dallas | STATUS_IN_PROGRESS | in    | 1         | 0         |
    And the following events are seeded for match "m-evts":
      | clock | typeID      | team          | players      |
      | 23'   | goal        | Columbus Crew | Hugo Picard  |
      | 39'   | yellow-card | FC Dallas     | Some Player  |
    And I am not logged in
    When I visit the matches page
    Then the now playing card for match "m-evts" should show event content
    And the now playing card for match "m-evts" should show "Hugo Picard"
    And the now playing card for match "m-evts" should show "Some Player"

  Scenario: Live match card hides events block when no events have occurred
    Given the following matches are seeded:
      | id      | homeTeam      | awayTeam  | status             | state | homeScore | awayScore |
      | m-noevt | Columbus Crew | FC Dallas | STATUS_IN_PROGRESS | in    | 0         | 0         |
    And I am not logged in
    When I visit the matches page
    Then the now playing card for match "m-noevt" should not show an events block

  Scenario: Live match card does not show substitutions on Now Playing
    Given the following matches are seeded:
      | id     | homeTeam      | awayTeam  | status             | state | homeScore | awayScore |
      | m-sub  | Columbus Crew | FC Dallas | STATUS_IN_PROGRESS | in    | 1         | 0         |
    And the following events are seeded for match "m-sub":
      | clock | typeID       | team          | players                       |
      | 23'   | goal         | Columbus Crew | Hugo Picard                   |
      | 60'   | substitution | Columbus Crew | Sub In Player, Sub Off Player |
    And I am not logged in
    When I visit the matches page
    Then the now playing card for match "m-sub" should show "Hugo Picard"
    And the now playing card for match "m-sub" should not show "Sub In Player"

  Scenario: Result card shows venue name
    Given the following matches are seeded:
      | id      | homeTeam      | awayTeam  | status           | homeScore | awayScore | venue                   |
      | m-ven-3 | Columbus Crew | FC Dallas | STATUS_FULL_TIME | 2         | 1         | ScottsMiracle-Gro Field |
    And the final score for match "m-ven-3" was 2-1 with Columbus home
    And I am not logged in
    When I visit the matches page
    Then the result card for match "m-ven-3" should show venue "ScottsMiracle-Gro Field"

  Scenario: Unpredicted upcoming match card shows team records and form
    Given the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status           | homeRecord | awayRecord | homeForm | awayForm |
      | m-form-1 | Columbus Crew | FC Dallas | STATUS_SCHEDULED | 5-3-2      | 4-4-2      | WWWLL    | LWDWL    |
    And I am not logged in
    When I visit the matches page
    Then the match card for "m-form-1" should show home record "5-3-2"
    And the match card for "m-form-1" should show home form "WWWLL"

  Scenario: Predicted upcoming match card shows team records and form
    Given the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status           | homeRecord | awayRecord | homeForm | awayForm |
      | m-form-2 | Columbus Crew | FC Dallas | STATUS_SCHEDULED | 5-3-2      | 4-4-2      | WWWLL    | LWDWL    |
    And I am logged in as "BlackAndGold@bsky.mock"
    And I have a seeded prediction of 2-1 for match "m-form-2"
    When I visit the matches page
    Then the match card for "m-form-2" should show home record "5-3-2"
    And the match card for "m-form-2" should show home form "WWWLL"

  Scenario: Now Playing match card shows team records and form
    Given the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status             | state | homeRecord | awayRecord | homeForm | awayForm |
      | m-form-3 | Columbus Crew | FC Dallas | STATUS_IN_PROGRESS | in    | 5-3-2      | 4-4-2      | WWWLL    | LWDWL    |
    And I am not logged in
    When I visit the matches page
    Then the now playing card for match "m-form-3" should show home record "5-3-2"
    And the now playing card for match "m-form-3" should show home form "WWWLL"

  Scenario: User's prediction appears at the top of a result card
    Given the following matches are seeded:
      | id       | homeTeam      | awayTeam  | status           | homeScore | awayScore |
      | m-pick-1 | Columbus Crew | FC Dallas | STATUS_FULL_TIME | 2         | 1         |
    And the final score for match "m-pick-1" was 2-1 with Columbus home
    And I am logged in as "BlackAndGold@bsky.mock"
    And I have a seeded prediction of 1-0 for match "m-pick-1"
    When I visit the matches page
    Then the result card for match "m-pick-1" should show my pick "1 – 0" below the score
