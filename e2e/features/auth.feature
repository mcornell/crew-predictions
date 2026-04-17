Feature: Authentication

  Scenario: Unauthenticated user sees sign-in option
    Given I am not logged in
    When I visit the matches page
    Then I should see a "Sign in with Google" link in the header

  Scenario: Sign in link points to the auth endpoint
    Given I am not logged in
    When I visit the matches page
    Then the "Sign in with Google" link should point to "/auth/login"

  Scenario: Logged-in user sees their name and a sign-out link
    Given I am logged in as "BlackYellow@bsky.social"
    When I visit the matches page
    Then I should see "BlackYellow@bsky.social" in the header
    And I should see a "Sign out" link in the header
    And I should not see a "Sign in with Google" link
