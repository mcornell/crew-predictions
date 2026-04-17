Feature: Authentication

  Scenario: Unauthenticated user sees sign-in option
    Given I am not logged in
    When I visit the matches page
    Then I should see a "Sign in with Google" link in the header

  Scenario: Sign in link points to the auth endpoint
    Given I am not logged in
    When I visit the matches page
    Then the "Sign in with Google" link should point to "/auth/login"
