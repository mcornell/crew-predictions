package org.acesradio.domain;

public record Prediction(
    Long id,
    Game game,
    User user,
    Integer homeGoals,
    Integer awayGoals
) {
    
}
