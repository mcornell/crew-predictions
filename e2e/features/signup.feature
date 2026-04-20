Feature: Sign up

  Scenario: New user can create an account with email and password
    When I visit the sign-up page
    And I sign up with email "newfan@example.com" and password "Nordecke96!"
    Then I should be on the matches page
    And I should see "newfan@example.com" in the header

  Scenario: Sign-up with an existing email shows an error
    Given a test user exists with email "existing@example.com" and password "Nordecke96!"
    When I visit the sign-up page
    And I sign up with email "existing@example.com" and password "DifferentPass1!"
    Then I should stay on the sign-up page
    And I should see the error "Could not create account"

  Scenario: Sign-up with a weak password shows an error
    When I visit the sign-up page
    And I sign up with email "weakpw@example.com" and password "abc"
    Then I should stay on the sign-up page
    And I should see the error "Could not create account"

  Scenario: Sign-up page links to login for existing users
    When I visit the sign-up page
    And I click the "Sign in" link
    Then I should be on the login page
