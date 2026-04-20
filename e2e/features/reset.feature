Feature: Password Reset

  Scenario: User can request a password reset email
    Given a test user exists with email "fan@example.com" and password "Nordecke96!"
    When I visit the reset page
    And I enter "fan@example.com" in the reset email field
    And I submit the reset form
    Then I should see a reset confirmation message

  Scenario: Login page has a forgot password link
    When I visit the login page
    Then I should see a "Forgot password?" link pointing to "/reset"
