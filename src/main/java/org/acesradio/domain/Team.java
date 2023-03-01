package org.acesradio.domain;

import io.micronaut.core.annotation.Nullable;
import io.micronaut.data.annotation.GeneratedValue;
import io.micronaut.data.annotation.Id;
import io.micronaut.data.annotation.MappedEntity;
import io.micronaut.serde.annotation.Serdeable;

import javax.validation.constraints.NotNull;

@Serdeable
@MappedEntity
public record Team(
        @Id @GeneratedValue @Nullable Long id,
        @NotNull String city,
        @NotNull String nickname
) {
}
