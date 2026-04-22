Feature: User profile

  Scenario: Profile page pre-populates the display name field with current handle
    Given I am logged in as "BlackAndGold@bsky.mock"
    When I visit the profile page
    Then the display name field should contain "BlackAndGold@bsky.mock"

  Scenario: Logged-in user can set a display name
    Given a test user exists with email "profiletest@crew.mock" and password "Nordecke96!"
    When I visit the login page
    And I sign in with email "profiletest@crew.mock" and password "Nordecke96!"
    Then I should be on the matches page
    When I visit the profile page
    And I set my display name to "Nordecke Regular"
    And I save my profile
    Then I should see "Nordecke Regular" in the header
    And I should not see "profiletest@crew.mock" in the header
