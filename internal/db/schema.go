package db

var schema = `
CREATE TABLE IF NOT EXISTS users (
    id serial primary key ,
    username varchar(40),
    telegram_id varchar(12)
);
`
