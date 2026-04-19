Feature: Login

  Scenario: User can sign in with email and password
    Given a test user exists with email "testfan@crew.mock" and password "Nordecke96!"
    When I visit the login page
    And I sign in with email "testfan@crew.mock" and password "Nordecke96!"
    Then I should be on the matches page
    And I should see "testfan@crew.mock" in the header
