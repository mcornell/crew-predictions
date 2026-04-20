Feature: Email verification banner

  Scenario: Unverified user sees a verification banner on the matches page
    Given a test user exists with email "unverified@crew.mock" and password "Nordecke96!"
    When I visit the login page
    And I sign in with email "unverified@crew.mock" and password "Nordecke96!"
    Then I should be on the matches page
    And I should see an email verification banner

  Scenario: Verified user does not see a verification banner
    Given I am logged in as a verified user
    When I visit the matches page
    Then I should not see an email verification banner
