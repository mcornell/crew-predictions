@reset
Feature: Mobile layout

  Background:
    Given the following matches are seeded:
      | id      | homeTeam      | awayTeam  | status           |
      | m-mob-1 | Columbus Crew | FC Dallas | STATUS_SCHEDULED |

  # ============================================================
  # Layout assertions per device viewport (consolidated — one
  # scenario per device, asserts everything that should hold on
  # that viewport using soft assertions for full failure visibility)
  # ============================================================

  Scenario: Matches page layout on iPhone 15
    Given I am not logged in
    And I am viewing on an iPhone 15
    When I visit the matches page
    Then the page should not overflow horizontally
    And I should see at least one Columbus Crew match card
    And each match card should be at most 260px tall
    And the site header should be at most 64px tall
    And the Predict button should be at least 44px tall

  Scenario: Matches page layout on Galaxy S24
    Given I am not logged in
    And I am viewing on a Galaxy S24
    When I visit the matches page
    Then the page should not overflow horizontally
    And I should see at least one Columbus Crew match card
    And each match card should be at most 260px tall
    And team names should not be clipped on any match card

  # ============================================================
  # Mobile auth flows (kept atomic — different code paths)
  # ============================================================

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

  # ============================================================
  # Hamburger menu — open, navigate, close (consolidated journey)
  # ============================================================

  Scenario: Hamburger menu opens, navigates, and closes
    Given I am not logged in
    And I am viewing on an iPhone 15
    When I visit the matches page
    And I tap the hamburger menu
    Then I should see the mobile navigation drawer
    When I tap outside the drawer
    Then the mobile navigation drawer should be closed
    When I tap the hamburger menu
    And I tap the Leaderboard link in the drawer
    Then I should be on the leaderboard page
