@reset
Feature: User profile

  Background:
    Given the following matches are seeded:
      | id        | homeTeam      | awayTeam  | status           |
      | match-p-1 | Columbus Crew | FC Dallas | STATUS_FULL_TIME |
      | match-p-2 | Columbus Crew | LA Galaxy | STATUS_FULL_TIME |

  Scenario: Profile page pre-populates the display name field with current handle
    Given I am logged in as "BlackAndGold@bsky.mock"
    When I visit my profile page
    Then the display name field should contain "BlackAndGold@bsky.mock"

  Scenario: Logged-in user can set a display name and location
    Given a test user exists with email "profiletest@crew.mock" and password "Nordecke96!"
    When I visit the login page
    And I sign in with email "profiletest@crew.mock" and password "Nordecke96!"
    Then I should be on the matches page
    When I visit my profile page
    And I set my display name to "Nordecke Regular"
    And I set my location to "Columbus, OH"
    And I save my profile
    Then I should see "Nordecke Regular" in the header
    And I should not see "profiletest@crew.mock" in the header

  Scenario: Profile page shows prediction count and leaderboard standing
    Given a test user exists with email "statsfan@crew.mock" and password "GoCrewGo1"
    And I visit the login page
    And I sign in with email "statsfan@crew.mock" and password "GoCrewGo1"
    And I should be on the matches page
    And I have a seeded prediction of 2-1 for match "match-p-1"
    And I have a seeded prediction of 1-0 for match "match-p-2"
    And the final score for match "match-p-1" was 2-1 with Columbus away
    And the final score for match "match-p-2" was 1-0 with Columbus away
    When I visit my profile page
    Then I should see my prediction count as 2
    And I should see my Aces Radio points

  Scenario: Location appears on user profile
    Given a test user exists with email "columbus@crew.mock" and password "GoCrewGo1"
    And I visit the login page
    And I sign in with email "columbus@crew.mock" and password "GoCrewGo1"
    And I should be on the matches page
    When I visit my profile page
    And I set my location to "Columbus, OH"
    And I save my profile
    And I visit my profile page
    Then the location field should contain "Columbus, OH"

  Scenario: Clicking a leaderboard handle navigates to that user's profile
    Given a test user exists with email "fanforprofile@crew.mock" and password "GoCrewGo1"
    And I visit the login page
    And I sign in with email "fanforprofile@crew.mock" and password "GoCrewGo1"
    And I should be on the matches page
    And I have a seeded prediction of 3-0 for match "match-p-1"
    And the final score for match "match-p-1" was 3-0 with Columbus away
    When I visit the leaderboard
    And I click the handle "fanforprofile@crew.mock" on the leaderboard
    Then I should be on the profile page for that user

  Scenario: Edit form only appears on own profile
    Given a test user exists with email "viewer@crew.mock" and password "GoCrewGo1"
    And a test user exists with email "other@crew.mock" and password "GoCrewGo1"
    And I visit the login page
    And I sign in with email "other@crew.mock" and password "GoCrewGo1"
    And I should be on the matches page
    And I have a seeded prediction of 2-1 for match "match-p-1"
    And I sign out
    And I visit the login page
    And I sign in with email "viewer@crew.mock" and password "GoCrewGo1"
    And I should be on the matches page
    And I have a seeded prediction of 1-0 for match "match-p-1"
    And the final score for match "match-p-1" was 1-0 with Columbus away
    When I visit the leaderboard
    And I click the handle "other@crew.mock" on the leaderboard
    Then I should not see the profile edit form
