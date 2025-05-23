-- Drop foreign key constraints first
ALTER TABLE users_observed_sensors
    DROP CONSTRAINT IF EXISTS "fk_users_observed_sensors_sensor_id_sensors_id";

ALTER TABLE users_observed_sensors
    DROP CONSTRAINT IF EXISTS "fk_users_observed_sensors_user_id_users_id";

ALTER TABLE users_observed_districts
    DROP CONSTRAINT IF EXISTS "fk_users_observed_districts_district_id_districts_id";

ALTER TABLE users_observed_districts
    DROP CONSTRAINT IF EXISTS "fk_users_observed_districts_user_id_users_id";

ALTER TABLE sensors
    DROP CONSTRAINT IF EXISTS "fk_sensors_district_id_districts_id";

-- Drop the created tables
DROP TABLE IF EXISTS users_observed_sensors;
DROP TABLE IF EXISTS users_observed_districts;

-- Remove the column added to users
ALTER TABLE users
    ALTER COLUMN operating_mode DROP DEFAULT;

ALTER TABLE users
    DROP COLUMN IF EXISTS operating_mode;