package org.aces.radio.domain;

    
public record Game(
    Long id,
    Team home,
    Team away,
    Short homeGoals,
    Short awayGoals) {
        
    }
