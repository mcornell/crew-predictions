package org.acesradio.domain;

import io.micronaut.core.annotation.Nullable;
import io.micronaut.data.annotation.GeneratedValue;
import io.micronaut.data.annotation.Id;
import io.micronaut.data.annotation.MappedEntity;
import io.micronaut.data.annotation.Relation;
import io.micronaut.serde.annotation.Serdeable;

import javax.validation.constraints.NotNull;

import static io.micronaut.data.annotation.Relation.Kind.MANY_TO_ONE;

@MappedEntity
@Serdeable
public record Prediction(
        @Id @GeneratedValue @Nullable Long id,
        @NotNull @Relation(MANY_TO_ONE) Handle handle,
        @NotNull @Relation(MANY_TO_ONE) Game game,
        @NotNull Short homeGoals,
        @NotNull Short awayGoals

) {
}

