Feature: Google sign-in

  Scenario: User can sign in with Google from the login page
    When I visit the login page
    And I sign in with Google as "gfan@example.com"
    Then I should be on the matches page
    And I should see "gfan@example.com" in the header

  Scenario: User can sign up with Google from the sign-up page
    When I visit the sign-up page
    And I sign in with Google as "newg@example.com"
    Then I should be on the matches page
    And I should see "newg@example.com" in the header

  Scenario: Google sign-in redirect works on iPhone 15
    Given I am viewing on an iPhone 15
    When I visit the login page
    And I sign in with Google as "gfan-ios@example.com"
    Then I should be on the matches page
    And I should see "gfan-ios@example.com" in the header

  Scenario: Google sign-in redirect works on Galaxy S24
    Given I am viewing on a Galaxy S24
    When I visit the login page
    And I sign in with Google as "gfan-android@example.com"
    Then I should be on the matches page
    And I should see "gfan-android@example.com" in the header
