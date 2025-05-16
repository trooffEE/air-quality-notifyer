ALTER TABLE sensors DROP COLUMN created_at;
ALTER TABLE sensors ADD timestamp TIMESTAMP;