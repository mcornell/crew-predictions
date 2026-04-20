Feature: Sign up

  Scenario: New user can create an account with email and password
    When I visit the sign-up page
    And I sign up with email "newfan@example.com" and password "Nordecke96!"
    Then I should be on the matches page
    And I should see "newfan@example.com" in the header
