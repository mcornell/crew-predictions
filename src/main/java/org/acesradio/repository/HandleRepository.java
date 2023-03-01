package org.acesradio.repository;

import io.micronaut.core.annotation.NonNull;
import io.micronaut.data.jdbc.annotation.JdbcRepository;
import io.micronaut.data.model.query.builder.sql.Dialect;
import io.micronaut.data.repository.PageableRepository;
import org.acesradio.domain.Handle;

import javax.validation.constraints.NotBlank;

@JdbcRepository(dialect = Dialect.POSTGRES)
public interface HandleRepository extends PageableRepository<Handle, Long> {

    Handle save(@NonNull @NotBlank String handle);
}

