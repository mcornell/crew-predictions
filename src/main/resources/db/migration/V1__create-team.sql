CREATE TABLE team(
    id        MEDIUMINT   not null,
    city      varchar(32) not null,
    nickname  varchar(32) not null,
    PRIMARY KEY (id)
)

INSERT INTO team (city, nickname) values
  ('Columbus', 'Crewzers SC'),
  ('Austin', 'Broccoli'),
  ('Los Angeles', 'Football Craft'),
  ('Philadelphia', 'Union');