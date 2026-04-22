@reset
Feature: User handles

  Background:
    Given the following matches are seeded:
      | id             | homeTeam      | awayTeam  | status           |
      | match-handle-1 | Columbus Crew | FC Dallas | STATUS_FULL_TIME |

  Scenario: Updated handle appears on leaderboard
    Given a test user exists with email "fan@crew.mock" and password "GoCrewGo1"
    And I visit the login page
    And I sign in with email "fan@crew.mock" and password "GoCrewGo1"
    And I should be on the matches page
    And I have a seeded prediction of 2-0 for match "match-handle-1"
    And the final score for match "match-handle-1" was 2-0 with Columbus away
    When I visit my profile page
    And I set my display name to "CrewForever"
    And I save my profile
    And I visit the leaderboard
    Then I should see "CrewForever" with 15 Aces Radio points
    And I should not see "fan@crew.mock" on the leaderboard
