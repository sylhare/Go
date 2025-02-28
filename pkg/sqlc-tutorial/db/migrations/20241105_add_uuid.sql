-- migrate:up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE authors
ADD COLUMN IF NOT EXISTS uuid UUID DEFAULT uuid_generate_v4();

-- migrate:down
ALTER TABLE authors DROP COLUMN IF EXISTS uuid;