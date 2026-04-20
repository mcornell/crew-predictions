Feature: Google sign-in

  Scenario: User can sign in with Google from the login page
    When I visit the login page
    And I sign in with Google as "gfan@example.com"
    Then I should be on the matches page
    And I should see "gfan@example.com" in the header
