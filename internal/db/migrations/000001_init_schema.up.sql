CREATE TABLE IF NOT EXISTS users (
    id serial primary key ,
    username varchar(40),
    telegram_id varchar(12)
);

CREATE TABLE IF NOT EXISTS districts (
    id serial primary key,
    name varchar(40),
    area geometry(multipolygon, 4326)
);