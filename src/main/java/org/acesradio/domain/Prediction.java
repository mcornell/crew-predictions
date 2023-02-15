package org.acesradio.domain;

public record Prediction(
    Long id,
    Game game,
    Handle handle,
    Integer homeGoals,
    Integer awayGoals
) {
    
}
