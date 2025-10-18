-- users_observed_districts
CREATE SEQUENCE users_observed_districts_id_seq;

ALTER TABLE users_observed_districts
    ALTER COLUMN id SET DEFAULT nextval('users_observed_districts_id_seq');

-- Set the sequence value to the maximum current id + 1
SELECT setval('users_observed_districts_id_seq', COALESCE(MAX(id), 0) + 1)
FROM users_observed_districts;

-- users_observed_sensors
CREATE SEQUENCE users_observed_sensors_id_seq;

ALTER TABLE users_observed_sensors
    ALTER COLUMN id SET DEFAULT nextval('users_observed_sensors_id_seq');

-- Set the sequence value to the maximum current id + 1
SELECT setval('users_observed_sensors_id_seq', COALESCE(MAX(id), 0) + 1)
FROM users_observed_sensors;