Feature: Prediction persistence

  Background:
    Given the following matches are seeded:
      | id              | homeTeam      | awayTeam           | status           |
      | match-persist-1 | Columbus Crew | Philadelphia Union | STATUS_SCHEDULED |

  Scenario: Submitted prediction survives a page reload
    Given I am logged in as "ColumbusNordecke@bsky.mock"
    When I visit the matches page
    And I enter a home score of 3 and away score of 1 for the first match
    And I click "Predict"
    And I revisit the matches page
    Then I should see my prediction of "3" on the first match card
