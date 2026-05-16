@reset
Feature: Match poll chain

  Scenario: Polling a single in-progress match records when it was polled
    Given the following matches are seeded:
      | id    | homeTeam      | awayTeam | status              | state |
      | m-301 | Columbus Crew | LAFC     | STATUS_IN_PROGRESS  | in    |
    When the admin triggers a score poll for match "m-301"
    Then the match detail for "m-301" reports a lastPollAt within the last 10 seconds
