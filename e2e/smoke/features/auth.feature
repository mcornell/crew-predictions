@smoke
Feature: Staging auth smoke tests

  Scenario: New user can sign up with email on desktop
    When I visit the staging sign-up page
    And I sign up with email "smoke-new@crew-predictions-staging.web.app"
    Then I should be on the staging matches page
    And I should see "smoke-new@crew-predictions-staging.web.app" in the staging header

  Scenario: Existing user can sign in with email on desktop
    When I visit the staging login page
    And I sign in with email "smoke-existing@crew-predictions-staging.web.app"
    Then I should be on the staging matches page
    And I should see "smoke-existing@crew-predictions-staging.web.app" in the staging header

  Scenario: New user can sign up with email on iPhone 15
    Given I am on an iPhone 15 viewport
    When I visit the staging sign-up page
    And I sign up with email "smoke-new@crew-predictions-staging.web.app"
    Then I should be on the staging matches page
    And I should see "smoke-new@crew-predictions-staging.web.app" in the staging header

  Scenario: Existing user can sign in with email on iPhone 15
    Given I am on an iPhone 15 viewport
    When I visit the staging login page
    And I sign in with email "smoke-existing@crew-predictions-staging.web.app"
    Then I should be on the staging matches page
    And I should see "smoke-existing@crew-predictions-staging.web.app" in the staging header

  Scenario: Existing user can sign in with email on Galaxy S24
    Given I am on a Galaxy S24 viewport
    When I visit the staging login page
    And I sign in with email "smoke-existing@crew-predictions-staging.web.app"
    Then I should be on the staging matches page
    And I should see "smoke-existing@crew-predictions-staging.web.app" in the staging header

  Scenario: Google sign-in initiates redirect on desktop
    When I visit the staging login page
    And I click the Google sign-in button
    Then the page should navigate toward Google for authentication

  Scenario: Google sign-in initiates redirect on iPhone 15
    Given I am on an iPhone 15 viewport
    When I visit the staging login page
    And I click the Google sign-in button
    Then the page should navigate toward Google for authentication
