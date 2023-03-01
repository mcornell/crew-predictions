package org.acesradio.domain;

import io.micronaut.core.annotation.Nullable;
import io.micronaut.data.annotation.GeneratedValue;
import io.micronaut.data.annotation.Id;
import io.micronaut.data.annotation.MappedEntity;
import io.micronaut.data.annotation.Relation;
import io.micronaut.serde.annotation.Serdeable;

import javax.validation.constraints.NotNull;

import java.util.List;

import static io.micronaut.data.annotation.Relation.Kind.ONE_TO_MANY;

@Serdeable
@MappedEntity
public record Handle(
        @Id @GeneratedValue @Nullable Long id,
        @NotNull String handle,
        @Relation(ONE_TO_MANY) @Nullable List<Prediction> predictions
) {
}