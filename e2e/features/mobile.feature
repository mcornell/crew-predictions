@reset
Feature: Mobile layout

  Background:
    Given the following matches are seeded:
      | id      | homeTeam      | awayTeam  | status           |
      | m-mob-1 | Columbus Crew | FC Dallas | STATUS_SCHEDULED |

  Scenario: Matches page fits within a mobile viewport without horizontal overflow
    Given I am not logged in
    And I am viewing on an iPhone 15
    When I visit the matches page
    Then the page should not overflow horizontally
    And I should see at least one Columbus Crew match card

  Scenario: Predict button meets minimum tap target size on mobile
    Given I am not logged in
    And I am viewing on an iPhone 15
    When I visit the matches page
    Then the Predict button should be at least 44px tall

  Scenario: Match cards do not collapse to an unreadable height on mobile
    Given I am not logged in
    And I am viewing on an iPhone 15
    When I visit the matches page
    Then each match card should be at most 180px tall

  Scenario: Site header stays within a single row on mobile
    Given I am not logged in
    And I am viewing on an iPhone 15
    When I visit the matches page
    Then the site header should be at most 64px tall

  Scenario: Matches page fits within a narrow Android viewport
    Given I am not logged in
    And I am viewing on a Galaxy S24
    When I visit the matches page
    Then the page should not overflow horizontally
    And I should see at least one Columbus Crew match card
    And each match card should be at most 180px tall

  Scenario: Team names are not clipped on Galaxy S24
    Given I am not logged in
    And I am viewing on a Galaxy S24
    When I visit the matches page
    Then team names should not be clipped on any match card

  Scenario: User can sign in with email and password on mobile
    Given a test user exists with email "mobilefan@example.com" and password "Nordecke96!"
    And I am viewing on an iPhone 15
    When I visit the login page
    And I sign in with email "mobilefan@example.com" and password "Nordecke96!"
    Then I should be on the matches page
    And I should see "mobilefan@example.com" in the header

  Scenario: User can sign in with Google on mobile
    Given I am viewing on an iPhone 15
    When I visit the login page
    And I sign in with Google as "gmobilefan@example.com"
    Then I should be on the matches page

  Scenario: Hamburger menu opens and closes on mobile
    Given I am not logged in
    And I am viewing on an iPhone 15
    When I visit the matches page
    And I tap the hamburger menu
    Then I should see the mobile navigation drawer
    When I tap outside the drawer
    Then the mobile navigation drawer should be closed

  Scenario: User can navigate via hamburger menu on mobile
    Given I am not logged in
    And I am viewing on an iPhone 15
    When I visit the matches page
    And I tap the hamburger menu
    And I tap the Leaderboard link in the drawer
    Then I should be on the leaderboard page
