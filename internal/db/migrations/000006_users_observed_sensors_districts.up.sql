ALTER TABLE users ADD COLUMN operating_mode SMALLINT;
ALTER TABLE users ALTER COLUMN operating_mode SET DEFAULT 0;

CREATE TABLE users_observed_districts (
    "id" integer NOT NULL,
    "user_id" integer NOT NULL,
    "district_id" integer NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE users_observed_sensors (
    "id" integer NOT NULL,
    "user_id" integer NOT NULL,
    "sensor_id" integer NOT NULL,
    PRIMARY KEY ("id")
);

ALTER TABLE sensors
    ADD CONSTRAINT "fk_sensors_district_id_districts_id" FOREIGN KEY("district_id") REFERENCES districts("id");

ALTER TABLE users_observed_districts
    ADD CONSTRAINT "fk_users_observed_districts_user_id_users_id" FOREIGN KEY("user_id") REFERENCES users("id");

ALTER TABLE users_observed_districts
    ADD CONSTRAINT "fk_users_observed_districts_district_id_districts_id" FOREIGN KEY("district_id") REFERENCES districts("id");

ALTER TABLE users_observed_sensors
    ADD CONSTRAINT "fk_users_observed_sensors_user_id_users_id" FOREIGN KEY("user_id") REFERENCES users("id");

ALTER TABLE users_observed_sensors
    ADD CONSTRAINT "fk_users_observed_sensors_sensor_id_sensors_id" FOREIGN KEY("sensor_id") REFERENCES sensors("id");