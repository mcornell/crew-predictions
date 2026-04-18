Feature: Authentication

  Scenario: Unauthenticated user sees sign-in option
    Given I am not logged in
    When I visit the matches page
    Then I should see a "Sign In" link in the header

  Scenario: Sign in link points to the login page
    Given I am not logged in
    When I visit the matches page
    Then the "Sign In" link should point to "/login"

  Scenario: Logged-in user sees their name and a sign-out link
    Given I am logged in as "BlackAndGold@bsky.mock"
    When I visit the matches page
    Then I should see "BlackAndGold@bsky.mock" in the header
    And I should see a "Sign out" link in the header
    And I should not see a "Sign In" link
