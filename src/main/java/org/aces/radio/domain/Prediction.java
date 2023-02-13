package org.aces.radio.domain;

public record Prediction(
    Long id,
    Game game,
    User user,
    Integer homeGoals,
    Integer awayGoals
) {
    
}
