package org.acesradio.repository;

import io.micronaut.data.jdbc.annotation.JdbcRepository;
import io.micronaut.data.repository.PageableRepository;

import io.micronaut.data.model.query.builder.sql.Dialect;
import org.acesradio.domain.Team;

@JdbcRepository(dialect = Dialect.POSTGRES)
public interface TeamRepository extends PageableRepository<Team, Long> {

}

