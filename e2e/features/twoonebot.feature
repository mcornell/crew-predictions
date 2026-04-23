@reset
Feature: TwoOneBot automated predictions

  Scenario: TwoOneBot predicts Columbus home match as 2-1
    Given the following matches are seeded:
      | id      | homeTeam      | awayTeam  | status           |
      | m-bot-1 | Columbus Crew | FC Dallas | STATUS_SCHEDULED |
    When the admin triggers a match refresh
    Then the predictions for match "m-bot-1" should include TwoOneBot with 2-1

  Scenario: TwoOneBot predicts Columbus away match as 1-2
    Given the following matches are seeded:
      | id      | homeTeam  | awayTeam      | status           |
      | m-bot-2 | FC Dallas | Columbus Crew | STATUS_SCHEDULED |
    When the admin triggers a match refresh
    Then the predictions for match "m-bot-2" should include TwoOneBot with 1-2

  Scenario: TwoOneBot does not predict a match already past kickoff
    Given the following matches are seeded in order:
      | id      | homeTeam      | awayTeam  | status           | kickoffOffset |
      | m-bot-3 | Columbus Crew | FC Dallas | STATUS_SCHEDULED | -1            |
    When the admin triggers a match refresh
    Then the predictions for match "m-bot-3" should not include TwoOneBot

  Scenario: TwoOneBot appears on the leaderboard after predicting
    Given the following matches are seeded:
      | id      | homeTeam      | awayTeam  | status           |
      | m-bot-4 | Columbus Crew | FC Dallas | STATUS_SCHEDULED |
    When the admin triggers a match refresh
    And I visit the leaderboard
    Then I should see "TwoOneBot" on the leaderboard
